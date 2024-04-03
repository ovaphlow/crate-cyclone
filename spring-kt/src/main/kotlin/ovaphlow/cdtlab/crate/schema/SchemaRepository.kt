package ovaphlow.cdtlab.crate.schema

import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository
import java.sql.ResultSet

@Repository
class SchemaRepository(private val jdbcTemplate: JdbcTemplate) {

    fun retrieve(): List<String> {
        val q = "select schema_name from information_schema.schemata"
        val result = jdbcTemplate.query(q) { rs: ResultSet, _: Int ->
            rs.getString("schema_name")
        }
        return result
    }

}
