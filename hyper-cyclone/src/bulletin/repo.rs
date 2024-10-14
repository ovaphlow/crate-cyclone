use chrono::{DateTime, Utc};
use tokio_postgres::Client;
use tokio_postgres::Error;

pub struct Bulletin {
    pub id: String,
    pub time: DateTime<Utc>,
    pub state: String,
    pub detail: String,
}

impl Bulletin {
    pub async fn update(client: &Client, id: String) -> Result<Self, Error> {
        let time = Utc::now();
        let bulletin = Bulletin {
            id: id,
            time: time,
            state: state,
            detail: "".to_string(),
        };
        client.execute(
            "UPDATE bulletin SET time = $1, state = $2 WHERE id = $3",
            &[&bulletin.time, &bulletin.state, &bulletin.id],
        ).await?;
        Ok(bulletin)
    }

    pub async fn retrieve_many(client: &Client) -> Result<Vec<Bulletin>, Error> {
        let mut bulletins = Vec::new();
        for row in client.query("SELECT id, time, state, detail FROM bulletin", &[]).await? {
            let bulletin = Bulletin {
                id: row.get(0),
                time: row.get(1),
                state: row.get(2),
                detail: row.get(3),
            };
            bulletins.push(bulletin);
        }
        Ok(bulletins)
    }

    pub async fn create(client: &Client, detail: String) -> Result<Self, Error> {
        let id = crate::utility::ksuid::generate_ksuid();
        let time = Utc::now();
        let bulletin = Bulletin {
            id: id,
            time: time,
            state: state,
            detail: detail,
        };

        client.execute(
            "INSERT INTO bulletin (id, time, state, detail) VALUES ($1, $2, $3, $4)",
            &[&bulletin.id, &bulletin.time, &bulletin.state, &bulletin.detail],
        ).await?;

        Ok(bulletin)
    }
}
