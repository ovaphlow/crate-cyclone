package ovaphlow.cdtlab.crate.table

import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository
import java.sql.ResultSet

@Repository
class TableRepository(private val jdbcTemplate: JdbcTemplate) {

    fun retrieve(schema: String): List<String> {
        val q = "select table_name from information_schema.tables where table_schema = ?"
        val result = jdbcTemplate.query(q, { rs: ResultSet, _: Int ->
            rs.getString("table_name")
        }, schema)
        return result
    }

}
