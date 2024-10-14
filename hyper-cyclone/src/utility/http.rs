use std::collections::HashMap;

use bytes::Bytes;
use http_body_util::{BodyExt, Full};

pub type BoxBody = http_body_util::combinators::BoxBody<Bytes, hyper::Error>;
type GenericError = Box<dyn std::error::Error + Send + Sync>;
pub type Result<T> = std::result::Result<T, GenericError>;

pub static STATUS_INTERNAL_SERVER_ERROR: &[u8] = b"Internal Server Error";
pub static STATUS_NOT_FOUND: &[u8] = b"Not Found";

pub fn full<T: Into<Bytes>>(chunk: T) -> BoxBody {
    Full::new(chunk.into())
        .map_err(|never| match never {})
        .boxed()
}

pub fn parse_query_string(query: &str) -> HashMap<String, String> {
    let mut params = HashMap::new();
    for pair in query.split('&') {
        let mut iter = pair.split('=');
        if let (Some(key), Some(value)) = (iter.next(), iter.next()) {
            params.insert(key.to_string(), value.to_string());
        }
    }
    params
}