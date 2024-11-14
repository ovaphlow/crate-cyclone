mod middleware;
mod router;
mod utility;

async fn hello(
    _: hyper::Request<hyper::body::Incoming>,
) -> Result<hyper::Response<http_body_util::Full<hyper::body::Bytes>>, std::convert::Infallible> {
    Ok(hyper::Response::new(http_body_util::Full::new(
        hyper::body::Bytes::from("crate cyclone with Rust and Hyper!"),
    )))
}

async fn router(
    req: hyper::Request<hyper::body::Incoming>,
) -> Result<hyper::Response<http_body_util::Full<hyper::body::Bytes>>, std::convert::Infallible> {
    match (req.method(), req.uri().path()) {
        (&hyper::Method::GET, "/") => hello(req).await,
        (&hyper::Method::GET, "/cyclone-api/bulletin") => router::bulletin::get(req).await,

        _ => {
            let mut not_found = hyper::Response::new(http_body_util::Full::new(
                hyper::body::Bytes::from("Not Found"),
            ));
            *not_found.status_mut() = hyper::StatusCode::NOT_FOUND;
            Ok(not_found)
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let addr = std::net::SocketAddr::from(([0, 0, 0, 0], 8422));

    let listener = tokio::net::TcpListener::bind(addr).await?;

    loop {
        let (stream, _) = listener.accept().await?;

        let io = hyper_util::rt::TokioIo::new(stream);

        tokio::task::spawn(async move {
            let svc = hyper::service::service_fn(router);
            let svc = tower::ServiceBuilder::new()
                .layer_fn(middleware::logger::Logger::new)
                .service(svc);
            if let Err(err) = hyper::server::conn::http1::Builder::new().serve_connection(io, svc).await {
                eprintln!("Error serving connection: {:?}", err);
            }
        });
    }
}
