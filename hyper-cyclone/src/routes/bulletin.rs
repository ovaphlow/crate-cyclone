use hyper::{header, Response, StatusCode};
use serde::{Deserialize, Serialize};

use crate::utilities::http::{BoxBody, Result, STATUS_INTERNAL_SERVER_ERROR};

#[derive(Serialize, Deserialize)]
struct Bulletin {
    id: u32,
    title: String,
    content: String,
}

impl Bulletin {
    fn new(id: u32, title: &str, content: &str) -> Self {
        Bulletin {
            id: id,
            title: title.to_string(),
            content: content.to_string(),
        }
    }
}

pub async fn handle_get(query: &str) -> Result<Response<BoxBody>> {
    let params = crate::utilities::http::parse_query_string(query);

    let id = params.get("id").and_then(|v| v.parse().ok()).unwrap_or(0);
    let default_string = "".to_string();
    let title = params.get("title").unwrap_or(&default_string);
    let content = params.get("content").unwrap_or(&default_string);

    let bulletin = Bulletin::new(id, title, content);
    let res = match serde_json::to_string(&bulletin) {
        Ok(json) => Response::builder()
            .header(header::CONTENT_TYPE, "application/json")
            .body(crate::utilities::http::full(json))
            .unwrap(),
        Err(_) => Response::builder()
            .status(StatusCode::INTERNAL_SERVER_ERROR)
            .body(crate::utilities::http::full(STATUS_INTERNAL_SERVER_ERROR))
            .unwrap(),
    };
    Ok(res)
}
