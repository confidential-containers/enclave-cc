use std::env;
use std::path::Path;
use std::sync::Arc;

use clap::{crate_authors, crate_version, App, Arg};
use protocols::{image, image_ttrpc};
use ttrpc::asynchronous::Server;
use ttrpc::{self, error::get_rpc_status as ttrpc_error};

use async_trait::async_trait;
use tokio::signal::unix::{signal, SignalKind};
use tokio::sync::Mutex;

use anyhow::{anyhow, Result};
use image_rs::config::ImageConfig;
use image_rs::image::ImageClient;
use image_rs::snapshots;

const CONTAINER_BASE: &str = "/run/enclave-cc-containers";
const CC_IMAGE_WORK_DIR: &str = "/run/image/";

// TODO: will replace with unix socket
const SOCK_ADDR: &str = "tcp://0.0.0.0:7788";

struct ImageService {
    image_client: Arc<Mutex<ImageClient>>,
}

// // This code was taken from https://github.com/kata-containers/kata-containers/blob/CCv0/src/agent/src/rpc.rs#L143 for quote
// A container ID must match this regex:
//     ^[a-zA-Z0-9][a-zA-Z0-9_.-]+$
fn verify_cid(id: &str) -> Result<()> {
    let mut chars = id.chars();

    let valid = match chars.next() {
        Some(first)
            if first.is_alphanumeric()
                && id.len() > 1
                && chars.all(|c| c.is_alphanumeric() || ['.', '-', '_'].contains(&c)) =>
        {
            true
        }
        _ => false,
    };

    match valid {
        true => Ok(()),
        false => Err(anyhow!("invalid container ID: {:?}", id)),
    }
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
                return Err(anyhow!("Invalid image name. {}", image));
            }
        } else {
            verify_cid(&cid)?;
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

//#[tokio::main(flavor = "current_thread")]
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
                .help("Specify the socket listen addr. Default is tcp://0.0.0.0:7788")
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

    println!("ttRPC server started");

    tokio::select! {
        _ = interrupt.recv() => {
            println!("shutdown the server");
            server.shutdown().await.unwrap();
        }
    };
}
