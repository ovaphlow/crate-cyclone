pub async fn handler_get_schema() -> &'static str {
    "hola el mundo!"
}
async fn repo_retrieve_schema_list() -> Result<Vec<String>, Box<dyn std::error::Error>> {
    let pool = crate::database::get_pool().await?;
    let client = pool.get().await?;
    let rows = client.query("select schema_name from information_schema.schemata", &[]).await?;
    let mut result = Vec::new();

    for row in rows {
        let value: String = row.get(0);
        result.push(value);
    }

    Ok(result)
}

async fn repo_retieve_table_list(schema: &str) -> Result<Vec<String>, Box<dyn std::error::Error>> {
    let pool = crate::database::get_pool().await?;
    let client = pool.get().await?;
    let rows = client.query("select table_name from information_schema.tables where table_schema = $1", &[&schema]).await?;
    let mut result = Vec::new();

    for row in rows {
        let value: String = row.get(0);
        result.push(value);
    }

    Ok(result)
}

async fn repo_retrieve_column_list(schema: &str, table: &str) -> Result<Vec<std::collections::HashMap<String, serde_json::Value>>, Box<dyn std::error::Error>> {
    let pool = crate::database::get_pool().await?;
    let client = pool.get().await?;
    let rows = client.query(r#"
        select ordinal_position, column_name, data_type
        from information_schema.columns
        where table_schema = $1 and table_name = $2
    "#, &[&schema, &table]).await?;
    let mut result = Vec::new();

    for row in rows {
        let mut map = std::collections::HashMap::new();
        map.insert("ordinal_position".to_string(), serde_json::Value::Number(serde_json::Number::from(row.get::<usize, i32>(0) as i64)));
        map.insert("column_name".to_string(), serde_json::Value::String(row.get(1)));
        map.insert("data_type".to_string(), serde_json::Value::String(row.get(2)));
        result.push(map);
    }

    Ok(result)
}
