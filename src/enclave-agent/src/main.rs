use std::sync::Arc;

use anyhow::Result;
use clap::Parser;
use log::{info, warn};
use protocols::image_ttrpc;
use tokio::signal::unix::{signal, SignalKind};
use ttrpc::asynchronous::Server;

mod config;
use config::{DecryptConfig, OcicryptConfig};

mod services;
use services::images::ImageService;

#[derive(Debug, Parser)]
#[command(version, long_about = None)]
struct Cli {
    /// enclave-agent socket listen address
    #[arg(default_value_t = String::from("tcp://0.0.0.0:7788"), short, long)]
    listen: String,

    /// decrypt configuration file path
    #[arg(default_value_t = String::default(), short = 'c', long)]
    decrypt_config: String,

    /// ocicrypt configuration file path
    #[arg(default_value_t = String::from("/etc/ocicrypt.conf"), short, long)]
    ocicrypt_config: String,
}

#[tokio::main(worker_threads = 1)]
async fn main() -> Result<()> {
    env_logger::Builder::from_env(env_logger::Env::default().default_filter_or("info")).init();

    let cli = Cli::parse();

    let dc = match DecryptConfig::load_from_file(&cli.decrypt_config) {
        Ok(config) => config,
        Err(e) => {
            warn!("Setting default DecryptConfig because: {:?}", e);
            DecryptConfig::default()
        }
    };

    OcicryptConfig::new(cli.ocicrypt_config)?.export_to_env();

    let image_service =
        Box::new(ImageService::new(dc)) as Box<dyn image_ttrpc::Image + Send + Sync>;

    let image_service = image_ttrpc::create_image(Arc::new(image_service));

    let mut server = Server::new()
        .bind(&cli.listen)?
        .register_service(image_service);

    let mut interrupt = signal(SignalKind::interrupt())?;
    server.start().await?;

    info!("ttRPC server started: {:?}", cli.listen);

    tokio::select! {
        _ = interrupt.recv() => {
            info!("shutdown the server");
            server.shutdown().await?;
        }
    };

    Ok(())
}
