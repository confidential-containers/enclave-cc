use std::path::Path;
use std::sync::Arc;

use anyhow::Result;
use async_trait::async_trait;
use image_rs::config::ImageConfig;
use image_rs::image::ImageClient;
use image_rs::snapshots;
use kata_sys_util::validate;
use log::info;
use protocols::image;
use tokio::sync::Mutex;
use ttrpc::{self, error::get_rpc_status as ttrpc_error};

use crate::config::DecryptConfig;

const CONTAINER_BASE: &str = "/run/enclave-cc/containers";

pub struct ImageService {
    dc: DecryptConfig,
    image_client: Arc<Mutex<ImageClient>>,
}

impl ImageService {
    pub fn new(dc: DecryptConfig) -> Self {
        let new_config = ImageConfig {
            default_snapshot: snapshots::SnapshotType::OcclumUnionfs,
            security_validate: dc.security_validate.map_or(true, |v| v),
            ..Default::default()
        };
        Self {
            dc,
            image_client: Arc::new(Mutex::new(ImageClient {
                config: new_config,
                ..Default::default()
            })),
        }
    }

    async fn pull_image(&self, req: &image::PullImageRequest) -> Result<String> {
        let image = req.get_image();
        let cid = self.get_container_id(req)?;
        let source_creds = (!req.get_source_creds().is_empty()).then(|| req.get_source_creds());
        let bundle_path = Path::new(CONTAINER_BASE).join(&cid);

        let dc_string = self
            .dc
            .key_provider
            .to_owned()
            .map_or(String::default(), |v| v);
        let dc_str = if dc_string.is_empty() {
            None
        } else {
            Some(dc_string.as_str())
        };

        info!("Pulling {:?}", image);
        self.image_client
            .lock()
            .await
            .pull_image(image, &bundle_path, &source_creds, &dc_str)
            .await?;

        Ok(image.to_owned())
    }

    fn get_container_id(&self, req: &image::PullImageRequest) -> Result<String> {
        let cid = req.get_container_id().to_string();
        // keep consistent with the kata container convention, more details
        // are described in https://github.com/confidential-containers/enclave-cc/issues/15
        validate::verify_id(&cid)?;
        Ok(cid)
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
                info!("Pull image {:?} successfully", r);
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

#[cfg(test)]
mod test {
    use super::*;

    #[rstest::rstest]
    #[case("redis", true)]
    #[case("redis_1.3", true)]
    #[case("redis:1.3", false)]
    #[case("", false)]
    fn test_get_container_id(#[case] container_id: &str, #[case] is_ok: bool) {
        let req = image::PullImageRequest {
            container_id: container_id.to_string(),
            ..Default::default()
        };

        let dc = DecryptConfig::load_from_file(
            &"test_data/decrypt_config/decrypt_config_normal.conf".to_string(),
        )
        .unwrap();
        let is = ImageService::new(dc);
        let res = is.get_container_id(&req);
        assert_eq!(res.is_ok(), is_ok, "err: {:?}", res);
    }
}
