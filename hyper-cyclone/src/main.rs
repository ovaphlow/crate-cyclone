use std::net::SocketAddr;

use hyper::{body::Incoming as IncomingBody, Method, Request, Response, StatusCode};
use hyper_util::rt::TokioIo;
use tokio::net::TcpListener;
use utility::http::{BoxBody, Result, STATUS_NOT_FOUND};

mod bulletin;
mod routes;
mod utility;

async fn handle_request(req: Request<IncomingBody>) -> Result<Response<BoxBody>> {
    match (req.method(), req.uri().path()) {
        (&Method::GET, "/crate-cyclone-api/bulletin") => {
            let query = req.uri().query().unwrap_or("");
            routes::bulletin::handle_get(query).await
        },
        _ => Ok(Response::builder()
            .status(StatusCode::NOT_FOUND)
            .body(crate::utility::http::full(STATUS_NOT_FOUND))
            .unwrap()),
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    let addr: SocketAddr = "127.0.0.1:8448".parse().unwrap();

    let listener = TcpListener::bind(&addr).await?;
    println!("Listening on http://{}", addr);

    loop {
        let (stream, _) = listener.accept().await?;
        let io = TokioIo::new(stream);

        tokio::task::spawn(async move {
            let service = hyper::service::service_fn(move |req| handle_request(req));

            if let Err(err) = hyper::server::conn::http1::Builder::new()
                .serve_connection(io, service)
                .await
            {
                eprintln!("Error serving connection: {:?}", err);
            }
        });
    }
}
