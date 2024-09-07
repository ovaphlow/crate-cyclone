use std::env;

pub async fn get_pool() -> Result<deadpool_postgres::Pool, tokio_postgres::Error> {
    dotenv::dotenv().ok();
    let user = env::var("DATABAASE_USER").expect("DATABAASE_USER 未设置");
    let password = env::var("DATABAASE_PASSWORD").expect("DATABAASE_PASSWORD 未设置");
    let host = env::var("DATABAASE_HOST").expect("DATABAASE_HOST 未设置");
    let port = env::var("DATABAASE_PORT").expect("DATABAASE_PORT 未设置").parse::<u16>().expect("数据库端口 错误");
    let name = env::var("DATABAASE_NAME").expect("DATABAASE_NAME 未设置");
    let mut cfg = deadpool_postgres::Config::new();
    cfg.user = Some(user);
    cfg.password = Some(password);
    cfg.host = Some(host);
    cfg.port = Some(port);
    cfg.dbname = Some(name);
    cfg.manager = Some(deadpool_postgres::ManagerConfig { recycling_method: deadpool_postgres::RecyclingMethod::Fast });
    let pool = cfg.create_pool(Some(deadpool_postgres::Runtime::Tokio1), tokio_postgres::NoTls).unwrap();
    Ok(pool)
}
