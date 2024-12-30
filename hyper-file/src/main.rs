use std::io::Read;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let addr = std::net::SocketAddr::from(([0, 0, 0, 0], 8422));

    let listener = tokio::net::TcpListener::bind(addr).await?;
    println!("Server listening on {}", addr);

    loop {
        let (stream, _) = listener.accept().await?;
        println!("Accepted connection from {:?}", stream.peer_addr());

        let io = hyper_util::rt::TokioIo::new(stream);

        tokio::task::spawn(async move {
            let svc = hyper::service::service_fn(router);
            let svc = tower::ServiceBuilder::new()
                .layer_fn(Logger::new)
                .service(svc);
            if let Err(err) = hyper::server::conn::http1::Builder::new()
                .serve_connection(io, svc)
                .await
            {
                eprintln!("Error serving connection: {:?}", err);
            }
        });
    }
}

/**
 * middleware logger
 */
#[derive(Debug, Clone)]
pub struct Logger<S> {
    inner: S,
}

impl<S> Logger<S> {
    pub fn new(inner: S) -> Self {
        Logger { inner }
    }
}

type Req = hyper::Request<hyper::body::Incoming>;

impl<S> hyper::service::Service<Req> for Logger<S>
where
    S: hyper::service::Service<Req>,
{
    type Response = S::Response;
    type Error = S::Error;
    type Future = S::Future;

    fn call(&self, req: Req) -> Self::Future {
        println!("Processing request: {} {}", req.method(), req.uri());
        self.inner.call(req)
    }
}

/**
 * router
 */
async fn router(
    req: hyper::Request<hyper::body::Incoming>,
) -> Result<hyper::Response<http_body_util::Full<hyper::body::Bytes>>, std::convert::Infallible> {
    match (req.method(), req.uri().path()) {
        (&hyper::Method::GET, "/crate-file-api/download") => get_file(req).await,

        _ => {
            let mut not_found = hyper::Response::new(http_body_util::Full::new(
                hyper::body::Bytes::from("Not Found"),
            ));
            *not_found.status_mut() = hyper::StatusCode::NOT_FOUND;
            Ok(not_found)
        }
    }
}

static USERS: &[(&str, &str, &str)] = &[
    ("user1", "password1", "/mnt/c/Users/ovaph/Desktop/squid_game_s02.jpg"),
    ("user2", "password2", "C:\\path\\to\\file2"),
];

// https://github.com/hyperium/hyper/blob/master/examples/send_file.rs
pub async fn get_file(
    req: hyper::Request<hyper::body::Incoming>,
) -> Result<hyper::Response<http_body_util::Full<hyper::body::Bytes>>, std::convert::Infallible> {
    let query = req.uri().query().unwrap_or("");
    let params: Vec<&str> = query.split('&').collect();
    let mut username = "";
    let mut password = "";

    for param in params {
        let kv: Vec<&str> = param.split('=').collect();
        if kv.len() == 2 {
            match kv[0] {
                "username" => username = kv[1],
                "password" => password = kv[1],
                _ => (),
            }
        }
    }

    for &(user, pass, file_path) in USERS {
        if user == username && pass == password {
            println!("Checking file path: {}", file_path); // 添加调试信息
            if std::path::Path::new(file_path).exists() {
                let mut file = std::fs::File::open(file_path).unwrap();
                let mut contents = Vec::new();
                file.read_to_end(&mut contents).unwrap();
                println!("File read successfully: {}", file_path);
                let content_length = contents.len();
                let response = hyper::Response::builder()
                    .status(hyper::StatusCode::OK)
                    .header("Content-Type", "application/octet-stream")
                    .header("Content-Length", content_length.to_string())
                    .body(http_body_util::Full::new(hyper::body::Bytes::from(contents)))
                    .unwrap();
                println!("Response built successfully with length: {}", content_length);
                return Ok(response);
            } else {
                println!("File not found: {}", file_path);
                return Ok(hyper::Response::builder()
                    .status(hyper::StatusCode::NOT_FOUND)
                    .body(http_body_util::Full::new(hyper::body::Bytes::from("文件不存在 (File not found)")))
                    .unwrap());
            }
        }
    }

    println!("Invalid username or password");
    Ok(hyper::Response::builder()
        .status(hyper::StatusCode::NOT_FOUND)
        .body(http_body_util::Full::new(hyper::body::Bytes::from("无效的账号或密码 (Invalid username or password)")))
        .unwrap())
}
