mod infrastructure;
mod routes;

use axum::{routing::get, Router};
use dotenv::dotenv;
use reqwest::Client;
use std::env;

#[tokio::main]
async fn main() {
    dotenv().ok();
    let port = env::var("HTTP_PORT").unwrap_or_else(|_| "8432".to_string());

    let app = Router::new()
        .route("/crate-api", get(|| async { "广告位 招租" }))
        .route("/crate-api/event", get(routes::event_api::hande_get_event));

    axum::Server::bind(&format!("0.0.0.0:{}", port).parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}
