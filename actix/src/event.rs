use std::collections::HashMap;

use actix_web::{web, HttpRequest, HttpResponse, Responder};
use serde::Serialize;
use serde_json::json;
use sqlx::prelude::*;

use crate::AppState;
use crate::condition_builder::array_contain_builder;
use crate::condition_builder::equal_builder;
use crate::condition_builder::greater_builder;
use crate::condition_builder::in_builder;
use crate::condition_builder::lesser_builder;
use crate::condition_builder::like_builder;
use crate::condition_builder::object_contain_builder;
use crate::condition_builder::object_like_builder;

const COLUMNS: [&str; 6] = ["id", "relation_id", "reference_id", "json_unquote(tags) tags", "json_unquote(detail) detail", "date_format(time, '%Y-%m-%d %H:%i:%s') time"];

#[derive(Serialize)]
struct Event {
    id: i64,
    #[serde(rename = "relationId")]
    relation_id: i64,
    #[serde(rename = "referenceId")]
    reference_id: i64,
    tags: String,
    detail: String,
    time: String,
    _id: String,
    #[serde(rename = "_relationId")]
    _relation_id: String,
    #[serde(rename = "_referenceId")]
    _reference_id: String,
}

pub async fn endpoint_event_get(app_data: web::Data<AppState>, req: HttpRequest) -> impl Responder {
    let query_string = req.query_string();
    let query: HashMap<String, String> = serde_qs::from_str(query_string).unwrap();
    if !query.contains_key("option") {
        return HttpResponse::NotAcceptable().json(json!({
            "error": "option is required"
        }));
    }
    if query.get("option") == Some(&"default".to_string()) {
        return filter_event_default(app_data, query).await;
    }
    HttpResponse::Ok().body("ok")
}

async fn filter_event_default(app_data: web::Data<AppState>, query: HashMap<String, String>) -> HttpResponse {
    let mut q = format!("select {} from events", COLUMNS.join(", "));
    let mut conditions: Vec<String> = Vec::new();
    let mut params: Vec<String> = Vec::new();
    if query.contains_key("equal") {
        let equal = query.get("equal").unwrap();
        let (c, p) = equal_builder(equal.split(",").collect::<Vec<&str>>().as_slice());
        conditions.extend(c);
        params.extend(p);
    }
    if query.contains_key("object-contain") {
        let object_contain = query.get("object-contain").unwrap();
        let (c, p) = object_contain_builder(object_contain.split(",").collect::<Vec<&str>>().as_slice());
        conditions.extend(c);
        params.extend(p);
    }
    if query.contains_key("array-contain") {
        let array_contain = query.get("array-contain").unwrap();
        let (c, p) = array_contain_builder(array_contain.split(",").collect::<Vec<&str>>().as_slice());
        conditions.extend(c);
        params.extend(p);
    }
    if query.contains_key("like") {
        let like = query.get("like").unwrap();
        let (c, p) = like_builder(like.split(",").collect::<Vec<&str>>().as_slice());
        conditions.extend(c);
        params.extend(p);
    }
    if query.contains_key("object-like") {
        let object_like = query.get("object-like").unwrap();
        let (c, p) = object_like_builder(object_like.split(",").collect::<Vec<&str>>().as_slice());
        conditions.extend(c);
        params.extend(p);
    }
    if query.contains_key("in") {
        let in_ = query.get("in").unwrap();
        let (c, p) = in_builder(in_.split(",").collect::<Vec<&str>>().as_slice());
        conditions.extend(c);
        params.extend(p);
    }
    if query.contains_key("lesser") {
        let lesser = query.get("lesser").unwrap();
        let (c, p) = lesser_builder(lesser.split(",").collect::<Vec<&str>>().as_slice());
        conditions.extend(c);
        params.extend(p);
    }
    if query.contains_key("greater") {
        let greater = query.get("greater").unwrap();
        let (c, p) = greater_builder(greater.split(",").collect::<Vec<&str>>().as_slice());
        conditions.extend(c);
        params.extend(p);
    }
    if conditions.len() > 0 {
        q.push_str(&format!(" where {}", conditions.join(" and ")));
    }
    q.push_str(&format!(" order by {} desc", "id"));
    q.push_str(&format!(" limit {}, {}", "0", "10"));
    let mut query = sqlx::query(&q);
    for param in &params {
        query = query.bind(param);
    }
    let result = query.fetch_all(&app_data.db).await;
    match result {
        Ok(rows) => {
            let events: Result<Vec<Event>, sqlx::Error> = rows
                .iter()
                .map(|row| {
                    Ok(Event {
                        id: row.get::<i64, _>("id"),
                        relation_id: row.get::<i64, _>("relation_id"),
                        reference_id: row.get::<i64, _>("reference_id"),
                        tags: row.get::<String, _>("tags"),
                        detail: row.get::<String, _>("detail"),
                        time: row.get::<String, _>("time"),
                        _id: row.get::<i64, _>("id").to_string(),
                        _relation_id: row.get::<i64, _>("relation_id").to_string(),
                        _reference_id: row.get::<i64, _>("reference_id").to_string(),
                    })
                })
                .collect();
            match events {
                Ok(events) => HttpResponse::Ok()
                    .append_header((crate::constants::HEADER_API_VERSION, "2024-01-06"))
                    .json(events),
                Err(e) => HttpResponse::InternalServerError().body(e.to_string()),
            }
        }
        Err(e) => HttpResponse::InternalServerError().body(e.to_string()),
    }
}
