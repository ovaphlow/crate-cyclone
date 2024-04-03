const COLUMNS: [&str; 6] = ["id", "email", "name", "phone", "json_unquote(tags) tags", "json_unquote(detail) detail", "date_format(time, '%Y-%m-%d %H:%i:%s') time"];

#[derive(Serialize)]
struct Subscriber {
    id: i64,
    email: String,
    name: String,
    phone: String,
    tags: String,
    detail: String,
    time: String,
    _id: String,
}

pub fn endpoint_sign_in(app_data: web::Data<AppState>, req: HttpRequest) -> impl Responder {
    HttpResponse::Ok()
        .append_header((crate::constants::HEADER_API_VERSION, "2024-01-06"))
        .body("ok")
}
