use anyhow::anyhow;
use serde::{Deserialize, Serialize};
use std::{fs::File, path::Path};

pub const DEFAULT_OCICRYPT_CONFIG_PATH: &str = "/etc/ocicrypt.conf";

#[derive(Serialize, Deserialize, Default)]
pub struct DecryptConfig {
    pub key_provider: Option<String>,
    #[serde(default = "default_security_validate")]
    pub security_validate: Option<bool>,
}

fn default_security_validate() -> Option<bool> {
    Some(true)
}

impl DecryptConfig {
    pub fn load_from_file(file_path: &String) -> Result<Self, anyhow::Error> {
        let file =
            File::open(file_path).map_err(|e| anyhow!("load decrypt config failed: {:?}", e))?;
        let config: DecryptConfig = serde_json::from_reader(file)
            .map_err(|e| anyhow!("parse decrypt config failed: {:?}", e))?;
        Ok(config)
    }
}

// OcicryptConfig is a config for the crate ocicrypt-rs.
pub struct OcicryptConfig {
    pub config_path: String,
}

impl OcicryptConfig {
    pub fn new(path: String) -> Result<Self, anyhow::Error> {
        if !Path::new(&path).exists() {
            return Err(anyhow!("ocicrypt config not found"));
        }
        Ok(Self { config_path: path })
    }
    pub fn export_to_env(&self) {
        std::env::set_var("OCICRYPT_KEYPROVIDER_CONFIG", self.config_path.clone());
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_load_decrypt_config() {
        struct ConfigCase {
            path: String,
            load_fail: bool,
            result_key_provider: Option<String>,
            result_security_validate: Option<bool>,
        }
        let decrypt_opt = "provider:attestation-agent:eaa_kbc::127.0.0.1:1234";
        let config_dir = "test_data/decrypt_config";
        let config_cases: Vec<ConfigCase> = vec![
            ConfigCase {
                path: format!("{}/{}", config_dir, "decrypt_config_not_existing.conf"),
                load_fail: true,
                result_key_provider: Some(decrypt_opt.to_string()),
                result_security_validate: Some(true),
            },
            ConfigCase {
                path: format!("{}/{}", config_dir, "decrypt_config_format_invalid.conf"),
                load_fail: true,
                result_key_provider: Some(decrypt_opt.to_string()),
                result_security_validate: Some(true),
            },
            ConfigCase {
                path: format!("{}/{}", config_dir, "decrypt_config_normal.conf"),
                load_fail: false,
                result_key_provider: Some(decrypt_opt.to_string()),
                result_security_validate: Some(true),
            },
            ConfigCase {
                path: format!("{}/{}", config_dir, "decrypt_config_disable_validate.conf"),
                load_fail: false,
                result_key_provider: Some(decrypt_opt.to_string()),
                result_security_validate: Some(false),
            },
            ConfigCase {
                path: format!("{}/{}", config_dir, "decrypt_config_default_validate.conf"),
                load_fail: false,
                result_key_provider: Some(decrypt_opt.to_string()),
                result_security_validate: Some(true),
            },
            ConfigCase {
                path: format!("{}/{}", config_dir, "decrypt_config_empty.conf"),
                load_fail: false,
                result_key_provider: None,
                result_security_validate: Some(true),
            },
        ];

        for case in config_cases {
            let result = DecryptConfig::load_from_file(&case.path);
            if case.load_fail {
                assert!(result.is_err());
            } else {
                let c = result.unwrap();
                assert_eq!(c.key_provider, case.result_key_provider);
                assert_eq!(c.security_validate, case.result_security_validate);
            }
        }
    }

    #[test]
    fn test_ocicrypt_config() {
        let config_path = "test_data/ocicrypt.conf".to_string();
        let oc = OcicryptConfig::new(config_path.clone()).unwrap();
        assert_eq!(oc.config_path, config_path);
        oc.export_to_env();
        assert_eq!(
            std::env::var("OCICRYPT_KEYPROVIDER_CONFIG").unwrap(),
            config_path
        );
    }
}
