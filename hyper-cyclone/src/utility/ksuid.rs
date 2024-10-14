use rand::Rng;
use std::time::{SystemTime, UNIX_EPOCH};

const EPOCH_OFFSET: u64 = 1_400_000_000; // Custom epoch offset for KSUID

pub fn generate_ksuid() -> String {
    // Get the current time since UNIX_EPOCH
    let now = SystemTime::now().duration_since(UNIX_EPOCH).expect("Time went backwards");
    let timestamp = now.as_secs() - EPOCH_OFFSET;

    // Generate 128 bits of random data
    let mut rng = rand::thread_rng();
    let mut payload = [0u8; 16];
    rng.fill(&mut payload);

    // Combine timestamp and payload
    let mut ksuid = Vec::with_capacity(20);
    ksuid.extend_from_slice(&timestamp.to_be_bytes()[4..]); // Use the last 4 bytes of the timestamp
    ksuid.extend_from_slice(&payload);

    // Encode as base62
    base62::encode(&ksuid)
}

fn main() {
    let ksuid = generate_ksuid();
    println!("Generated KSUID: {}", ksuid);
}

mod base62 {
    const CHARSET: &[u8] = b"0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";

    pub fn encode(data: &[u8]) -> String {
        let mut value = 0u128;
        for &byte in data {
            value = (value << 8) | byte as u128;
        }

        let mut encoded = Vec::new();
        while value > 0 {
            let remainder = (value % 62) as usize;
            value /= 62;
            encoded.push(CHARSET[remainder]);
        }

        encoded.reverse();
        String::from_utf8(encoded).expect("Invalid UTF-8")
    }
}