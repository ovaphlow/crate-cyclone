use tokio_postgres::Client;

use super::repo::Bulletin;

pub async fn save_bulletin(client: &Client, detail: String) -> Result<Bulletin, tokio_postgres::Error> {
    let bulletin = Bulletin::create(client, detail).await?;
    println!("Bulletin saved: {:?}", bulletin);
    Ok(bulletin)
}