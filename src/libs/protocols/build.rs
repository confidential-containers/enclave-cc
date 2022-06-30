use std::path::Path;
use ttrpc_codegen::{Codegen, Customize};

fn main() {
    let protos = vec!["protos/image.proto"];

    // Tell Cargo that if the .proto files changed, to rerun this build script.
    protos
        .iter()
        .for_each(|p| println!("cargo:rerun-if-changed={}", &p));

    let out_dir = Path::new("src");

    Codegen::new()
        .out_dir(out_dir)
        .inputs(&protos)
        .include("protos")
        .rust_protobuf()
        .customize(Customize {
            async_all: true,
            ..Default::default()
        })
        .run()
        .expect("Gen code failed.");
}
