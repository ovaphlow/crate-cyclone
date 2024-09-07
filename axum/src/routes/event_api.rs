use axum::{response::IntoResponse, Json};
use serde_json::json;

pub async fn hande_get_event() -> impl IntoResponse {
    Json(json!({
        "id": 1,
        "name": "Rust 2021",
        "location": "Online",
        "start_date": "2021-05-15",
        "end_date": "2021-05-16",
    }))
}
