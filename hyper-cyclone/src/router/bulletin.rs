pub async fn bulletin(
    _: hyper::Request<hyper::body::Incoming>,
) -> Result<hyper::Response<http_body_util::Full<hyper::body::Bytes>>, std::convert::Infallible> {
    Ok(hyper::Response::new(http_body_util::Full::new(
        hyper::body::Bytes::from("Hello, World! bulletin"),
    )))
}
