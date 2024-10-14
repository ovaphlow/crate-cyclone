use tokio_postgres::{Client, Error, NoTls};

pub async fn init_post_connection() -> Result<Client, Error> {
    let (client, connection) = 
        tokio_postgres::connect("host=localhost user=ovaphlow dbname=postgres password=", NoTls).await?;
    
    tokio::spawn(async move{
        if let Err(e) = connection.await {
            eprintln!("connection error: {}", e);
        }
    });

    Ok(client)
}