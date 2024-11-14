use hyper::server::conn::http1;
use hyper::service::service_fn;

mod router;
mod utility;

use router::bulletin::{self, bulletin};

async fn hello(
    _: hyper::Request<hyper::body::Incoming>,
) -> Result<hyper::Response<http_body_util::Full<hyper::body::Bytes>>, std::convert::Infallible> {
    Ok(hyper::Response::new(http_body_util::Full::new(
        hyper::body::Bytes::from("Hello, World!"),
    )))
}

async fn hello1(
    _: hyper::Request<hyper::body::Incoming>,
) -> Result<hyper::Response<http_body_util::Full<hyper::body::Bytes>>, std::convert::Infallible> {
    Ok(hyper::Response::new(http_body_util::Full::new(
        hyper::body::Bytes::from("Hello, World!111"),
    )))
}

async fn router(
    req: hyper::Request<hyper::body::Incoming>,
) -> Result<hyper::Response<http_body_util::Full<hyper::body::Bytes>>, std::convert::Infallible> {
    match (req.method(), req.uri().path()) {
        (&hyper::Method::GET, "/hello") => hello(req).await,
        (&hyper::Method::GET, "/hello1") => hello1(req).await,
        (&hyper::Method::GET, "/cyclone-api/bulletin") => bulletin(req).await,

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
    let addr = std::net::SocketAddr::from(([127, 0, 0, 1], 8422));

    let listener = tokio::net::TcpListener::bind(addr).await?;

    loop {
        let (stream, _) = listener.accept().await?;

        let io = hyper_util::rt::TokioIo::new(stream);

        tokio::task::spawn(async move {
            if let Err(err) = http1::Builder::new()
                .serve_connection(io, service_fn(router))
                .await
            {
                eprintln!("Error serving connection: {:?}", err);
            }
        });
    }
}
