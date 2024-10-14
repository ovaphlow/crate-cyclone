use hyper::{header, Response, StatusCode};
use serde::{Deserialize, Serialize};
use serde_json::json;

use crate::utility::http::{BoxBody, Result, STATUS_INTERNAL_SERVER_ERROR};

pub async fn handle_get(query: &str) -> Result<Response<BoxBody>, hyper::Error> {
    let params = crate::utility::http::parse_query_string(query);

    let detail = params.get("detail").unwrap_or(&"".to_string());

    match crate::bulletin::service::save_bulletin(client, detail.to_string()).await {
        Ok(bulletin) => {
            let json = json!({
                "id": bulletin.id,
                "time": bulletin.time,
                "state": bulletin.state,
                "detail": bulletin.detail,
            });

            let response = Response::builder()
                .header("Content-Type", "application/json")
                .body(Body::from(json.to_string()))
                .unwrap();

            Ok(response)
        }
        Err(e) => {
            eprintln!("Error: {:?}", e);
            let res = Response::builder()
                .status(StatusCode::INTERNAL_SERVER_ERROR)
                .body(crate::utility::http::full(STATUS_INTERNAL_SERVER_ERROR))
                .unwrap();
            Ok(res)
        }
    }
}
