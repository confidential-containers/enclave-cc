use std::env;
use std::fs;
use std::fs::File;
use std::io::prelude::*;
use std::path::Path;
use std::sync::Arc;

use anyhow::{anyhow, Result};
use async_trait::async_trait;
use clap::{crate_authors, crate_version, App, Arg};
use image_rs::config::ImageConfig;
use image_rs::image::ImageClient;
use image_rs::snapshots;
use kata_sys_util::validate;
use protocols::{image, image_ttrpc};
use tokio::signal::unix::{signal, SignalKind};
use tokio::sync::Mutex;
use ttrpc::asynchronous::Server;
use ttrpc::{self, error::get_rpc_status as ttrpc_error};

const CONTAINER_BASE: &str = "/run/enclave-cc/containers";

// TODO: will replace with unix socket
const TCP_SOCK_ADDR: &str = "tcp://0.0.0.0:7788";

struct ImageService {
    image_client: Arc<Mutex<ImageClient>>,
}

impl ImageService {
    fn new() -> Self {
        let new_config = ImageConfig {
            default_snapshot: snapshots::SnapshotType::OcclumUnionfs,
            ..Default::default()
        };
        Self {
            image_client: Arc::new(Mutex::new(ImageClient {
                config: new_config,
                ..Default::default()
            })),
        }
    }

    async fn pull_image(&self, req: &image::PullImageRequest) -> Result<String> {
        let image = req.get_image();
        let mut cid = req.get_container_id().to_string();

        if cid.is_empty() {
            let v: Vec<&str> = image.rsplit('/').collect();
            if !v[0].is_empty() {
                // ':' have special meaning for umoci during upack
                cid = v[0].replace(":", "_");
            } else {
                return Err(anyhow!("Invalid image name: {:?}", image));
            }
        } else {
            validate::verify_id(&cid)?;
        }

        let keyprovider_config = Path::new("/etc").join("ocicrypt_keyprovider_native.conf");
        if !keyprovider_config.exists() {
            let config = r#"
            {
                "key-providers": {
                    "attestation-agent": {
                        "native": "attestation-agent"
                    }
                }
            }
            "#;
            File::create(&keyprovider_config)?.write_all(config.as_bytes())?;
        }
        std::env::set_var("OCICRYPT_KEYPROVIDER_CONFIG", keyprovider_config);

        let mut config: String = String::new();
        let args: Vec<String> = env::args().collect();
        // If config file specified in the args, read contents from config file
        let config_position = args.iter().position(|a| a == "--decrypt-config" || a == "-c");
        if let Some(config_position) = config_position {
            if let Some(config_file) = args.get(config_position + 1) {
                config = fs::read_to_string(config_file).expect("Config file not found");
            } else {
                panic!("The config argument wasn't formed properly: {:?}", args);
            }
        }

        let decrypt_config = if !config.is_empty() {
            Some(config.as_str())
        } else {
            None
        };

        let source_creds = (!req.get_source_creds().is_empty()).then(|| req.get_source_creds());

        let bundle_path = Path::new(CONTAINER_BASE).join(&cid);

        println!("Pulling {:?}", image);
        self.image_client
            .lock()
            .await
            .pull_image(image, &bundle_path, &source_creds, &decrypt_config)
            .await?;

        Ok(image.to_owned())
    }
}

#[async_trait]
impl protocols::image_ttrpc::Image for ImageService {
    async fn pull_image(
        &self,
        _ctx: &ttrpc::r#async::TtrpcContext,
        req: image::PullImageRequest,
    ) -> ttrpc::Result<image::PullImageResponse> {
        match self.pull_image(&req).await {
            Ok(r) => {
                println!("Pull image {:?} successfully", r.clone());
                let mut resp = image::PullImageResponse::new();
                resp.image_ref = r;
                return Ok(resp);
            }
            Err(e) => {
                return Err(ttrpc_error(ttrpc::Code::INTERNAL, e.to_string()));
            }
        }
    }
}

#[tokio::main(worker_threads = 1)]
async fn main() -> Result<()> {
    let matches = App::new("Enclave agent")
        .version(crate_version!())
        .author(crate_authors!())
        .arg(
            Arg::with_name("listen")
                .short("l")
                .long("listen")
                .value_name("sockaddr")
                .help(&format!(
                    "{}{}",
                    "Specify the socket listen addr. Default is ", TCP_SOCK_ADDR
                ))
                .takes_value(true),
        )
        .arg(
            Arg::with_name("decrypt-config")
                .short("c")
                .long("decrypt-config")
                .help(&format!("The decrypt config file path"))
                .takes_value(true)
                .required(false),
        )
        .get_matches();

    let sockaddr = if let Some(addr) = matches.value_of("listen") {
        addr
    } else {
        TCP_SOCK_ADDR
    };

    let image_service = Box::new(ImageService::new()) as Box<dyn image_ttrpc::Image + Send + Sync>;

    let image_service = image_ttrpc::create_image(Arc::new(image_service));

    let mut server = Server::new()
        .bind(sockaddr)?
        .register_service(image_service);

    let mut interrupt = signal(SignalKind::interrupt())?;
    server.start().await?;

    println!("ttRPC server started: {:?}", sockaddr);

    tokio::select! {
        _ = interrupt.recv() => {
            println!("shutdown the server");
            server.shutdown().await?;
        }
    };

    Ok(())
}
