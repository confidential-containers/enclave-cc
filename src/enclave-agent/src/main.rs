use std::env;
use std::sync::Arc;

use anyhow::Result;
use clap::{crate_authors, crate_version, App, Arg};
use protocols::image_ttrpc;
use tokio::signal::unix::{signal, SignalKind};
use ttrpc::asynchronous::Server;

mod config;
use config::{DecryptConfig, OcicryptConfig, DEFAULT_OCICRYPT_CONFIG_PATH};

mod services;
use services::images::ImageService;

// TODO: will replace with unix socket
const TCP_SOCK_ADDR: &str = "tcp://0.0.0.0:7788";

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
                .help("The decrypt config file path".as_ref())
                .takes_value(true)
                .required(false),
        )
        .arg(
            Arg::with_name("ocicrypt-config")
                .short("o")
                .long("ocicrypt-config")
                .help("The ocicrypt config file path".as_ref())
                .takes_value(true)
                .required(false),
        )
        .get_matches();

    let sockaddr = if let Some(addr) = matches.value_of("listen") {
        addr
    } else {
        TCP_SOCK_ADDR
    };

    let dc = if let Some(file_path) = matches.value_of("decrypt-config") {
        DecryptConfig::load_from_file(&file_path.to_string())?
    } else {
        DecryptConfig::default()
    };

    let oc = if let Some(file_path) = matches.value_of("ocicrypt-config") {
        OcicryptConfig::new(file_path.to_string())?
    } else {
        OcicryptConfig::new(DEFAULT_OCICRYPT_CONFIG_PATH.to_string())?
    };
    oc.export_to_env();

    let image_service =
        Box::new(ImageService::new(dc)) as Box<dyn image_ttrpc::Image + Send + Sync>;

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
