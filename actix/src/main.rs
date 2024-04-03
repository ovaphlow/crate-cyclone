mod condition_builder;
mod constants;
mod event;
mod subscriber;

use actix_web::{web, App, HttpRequest, HttpResponse, HttpServer, Responder};
use dotenv::dotenv;
use sqlx::Pool;
use sqlx::mysql::MySqlSslMode;
use sqlx::mysql::MySqlConnectOptions;
use sqlx::mysql::MySqlPoolOptions;
use std::env;

pub struct AppState {
    db: Pool<sqlx::MySql>,
}

async fn index(req: HttpRequest) -> impl Responder {
    print!("{} {}\n", req.method().to_string(), req.uri().to_string());
    HttpResponse::Ok().body("Hello world")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();

    let database_username: String = env::var("DATABASE_USERNAME").expect("DATABASE_USERNAME must be set");
    let database_password: String = env::var("DATABASE_PASSWORD").expect("DATABASE_PASSWORD must be set");
    let database_host: String = env::var("DATABASE_HOST").expect("DATABASE_HOST must be set");
    let database_port_str: String = env::var("DATABASE_PORT").expect("DATABASE_PORT must be set");
    let database_port: u16 = database_port_str.parse().expect("DATABASE_PORT must be a number");
    let database_name: String = env::var("DATABASE_NAME").expect("DATABASE_NAME must be set");
    let options: MySqlConnectOptions = MySqlConnectOptions::new()
        .username(&database_username)
        .password(&database_password)
        .host(&database_host)
        .port(database_port)
        .database(&database_name)
        .ssl_mode(MySqlSslMode::Disabled);
    let pool: Pool<sqlx::MySql> = match MySqlPoolOptions::new()
        .max_connections(5)
        .connect_with(options)
        .await
    {
        Ok(pool) => {
            println!("✅Connection to the database is successful!");
            pool
        }
        Err(err) => {
            println!("🈚Failed to connect to the database: {:?}", err);
            std::process::exit(1);
        },
    };
    println!("Connected to {}", database_host);

    let server_address: String = env::var("SERVER_ADDRESS").expect("SERVER_ADDRESS must be set");
    let server_port: u16 = env::var("SERVER_PORT")
        .expect("SERVER_PORT must be set")
        .parse()
        .expect("SERVER_PORT must be a number");

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(AppState { db: pool.clone() }))
            .route("/", web::get().to(index))
            .route("/crate-api/event", web::get().to(event::endpoint_event_get))
    })
    .bind((server_address, server_port))?
    .run()
    .await
}
