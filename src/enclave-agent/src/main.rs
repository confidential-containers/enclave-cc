use std::env;
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
const CC_IMAGE_WORK_DIR: &str = "/run/image/";

// TODO: will replace with unix socket
const SOCK_ADDR: &str = "tcp://0.0.0.0:7788";

struct ImageService {
    image_client: Arc<Mutex<ImageClient>>,
}

impl ImageService {
    fn new() -> Self {
        env::set_var("CC_IMAGE_WORK_DIR", CC_IMAGE_WORK_DIR);
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

        let source_creds = (!req.get_source_creds().is_empty()).then(|| req.get_source_creds());

        let bundle_path = Path::new(CONTAINER_BASE).join(&cid);

        self.image_client
            .lock()
            .await
            .pull_image(image, &bundle_path, &source_creds, &None)
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
async fn main() {
    let matches = App::new("Enclave agent")
        .version(crate_version!())
        .author(crate_authors!())
        .arg(
            Arg::with_name("listen")
                .short("l")
                .long("listen")
                .value_name("sockaddr")
                .help(
                    ("Specify the socket listen addr. Default is ".to_string()
                        + format!("{}", SOCK_ADDR).as_str())
                    .as_str(),
                )
                .takes_value(true),
        )
        .get_matches();

    let mut sockaddr = SOCK_ADDR;
    if matches.is_present("listen") {
        sockaddr = matches.value_of("listen").unwrap();
    }

    let image_service = Box::new(ImageService::new()) as Box<dyn image_ttrpc::Image + Send + Sync>;

    let image_service = image_ttrpc::create_image(Arc::new(image_service));

    let mut server = Server::new()
        .bind(sockaddr)
        .unwrap()
        .register_service(image_service);

    let mut interrupt = signal(SignalKind::interrupt()).unwrap();
    server.start().await.unwrap();

    println!("ttRPC server started: {:?}", sockaddr);

    tokio::select! {
        _ = interrupt.recv() => {
            println!("shutdown the server");
            server.shutdown().await.unwrap();
        }
    };
}
