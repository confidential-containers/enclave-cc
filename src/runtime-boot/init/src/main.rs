extern crate libc;

use libc::syscall;

use nix::mount::MsFlags;
use std::env;
use std::error::Error;
use std::fs::File;
use std::io::{ErrorKind, Read};

use anyhow::Result;
use std::ffi::CString;
use std::mem::size_of;
use std::path::Path;

fn main() -> Result<(), Box<dyn Error>> {
    let agent_boot = matches!(env::var("ENCLAVE_AGENT"), Ok(val) if val == "true" || val == "TRUE" || val == "1");

    // Mount the image
    const SYS_MOUNT_FS: i64 = 363;
    const KEY_FILE: &str = "/tmp/key.txt";

    let ret = match agent_boot {
        true => {
            let null_ptr: *const i8 = std::ptr::null();
            unsafe { syscall(SYS_MOUNT_FS, null_ptr, null_ptr) }
        }
        false => {
            let rootfs_upper_layer = "/sefs/upper";
            let rootfs_lower_layer = "/sefs/lower";
            let rootfs_entry = "/";

            let fs_type = String::from("sefs");
            let source = Path::new(&fs_type);

            let mount_path = Path::new("/tmp");
            let flags = MsFlags::empty();

            nix::mount::mount(
                Some(source),
                mount_path,
                Some(fs_type.as_str()),
                flags,
                Some("dir=/keys/sefs/lower"),
            )
            .unwrap_or_else(|err| {
                eprintln!("Error mounting keys: {}", err);
            });

            // Get the key of FS image
            let key = {
                let key_str = load_key(KEY_FILE)?;
                let mut key: sgx_key_128bit_t = Default::default();
                parse_str_to_bytes(&key_str, &mut key)?;
                Some(key)
            };
            nix::mount::umount(mount_path)?;
            let key_ptr = key
                .as_ref()
                .map(|key| key as *const sgx_key_128bit_t)
                .unwrap_or(std::ptr::null());

            // Example envs. must end with null
            let env1 = CString::new("TEST=1234").unwrap();
            let envp = [env1.as_ptr(), std::ptr::null()];
            // Set rootfs parameters
            let upper_layer_path = CString::new(rootfs_upper_layer).expect("CString::new failed");
            let lower_layer_path = CString::new(rootfs_lower_layer).expect("CString::new failed");
            let entry_point = CString::new(rootfs_entry).expect("CString::new failed");
            let hostfs_source = CString::new("/tmp").expect("CString::new failed");
            let rootfs_config: user_rootfs_config = user_rootfs_config {
                len: size_of::<user_rootfs_config>(),
                upper_layer_path: upper_layer_path.as_ptr(),
                lower_layer_path: lower_layer_path.as_ptr(),
                entry_point: entry_point.as_ptr(),
                hostfs_source: hostfs_source.as_ptr(),
                hostfs_target: std::ptr::null(),
                envp: envp.as_ptr(),
            };
            unsafe { syscall(SYS_MOUNT_FS, key_ptr, &rootfs_config) }
        }
    };
    if ret < 0 {
        return Err(Box::new(std::io::Error::last_os_error()));
    }
    Ok(())
}

#[allow(non_camel_case_types)]
type sgx_key_128bit_t = [u8; 16];

#[repr(C)]
#[derive(Debug, Copy, Clone)]
#[allow(non_camel_case_types)]
struct user_rootfs_config {
    // length of the struct
    len: usize,
    // UnionFS type rootfs upper layer, read-write layer
    upper_layer_path: *const i8,
    // UnionFS type rootfs lower layer, read-only layer
    lower_layer_path: *const i8,
    entry_point: *const i8,
    // HostFS source path
    hostfs_source: *const i8,
    // HostFS target path, default value is "/host"
    hostfs_target: *const i8,
    // An array of pointers to null-terminated strings
    // and must be terminated by a null pointer
    envp: *const *const i8,
}

fn load_key(key_path: &str) -> Result<String, Box<dyn Error>> {
    let mut key_file = File::open(key_path)?;
    let mut key = String::new();
    key_file.read_to_string(&mut key)?;
    Ok(key.trim_end_matches(|c| c == '\r' || c == '\n').to_string())
}

fn parse_str_to_bytes(arg_str: &str, bytes: &mut [u8]) -> Result<(), Box<dyn Error>> {
    let bytes_str_vec = {
        let bytes_str_vec: Vec<&str> = arg_str.split('-').collect();
        if bytes_str_vec.len() != bytes.len() {
            return Err(Box::new(std::io::Error::new(
                ErrorKind::InvalidData,
                "The length or format of Key/MAC string is invalid",
            )));
        }
        bytes_str_vec
    };

    for (byte_i, byte_str) in bytes_str_vec.iter().enumerate() {
        bytes[byte_i] = u8::from_str_radix(byte_str, 16)?;
    }
    Ok(())
}
