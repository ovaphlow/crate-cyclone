mod database;
mod database_schema;
extern crate snowflake;

use axum::{
    http::StatusCode,
    response::{IntoResponse, Json},
    routing::{get},
    Router,
};
use serde_json::json;
use snowflake::ProcessUniqueId;

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/crate-api", get(|| async { "广告位 招租"}))
        .route("/crate-api/schema", get(table::handler_schema_get));

    axum::Server::bind(&"0.0.0.0:8421".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
