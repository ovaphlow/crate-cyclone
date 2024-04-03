package ovaphlow.cdtlab.crate.generalimpl

import com.fasterxml.jackson.databind.ObjectMapper
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository
import ovaphlow.cdtlab.crate.utilities.IdGenerator
import ovaphlow.cdtlab.crate.utility.ConditionBuilder
import java.sql.ResultSet
import java.util.*

@Repository
class Repository(private val jdbcTemplate: JdbcTemplate) {

    fun retrieveColumns(schema: String, table: String): List<Column> {
        val q = "select column_name from information_schema.columns where table_schema = ? and table_name = ?"
        val result = jdbcTemplate.query(q, { rs: ResultSet, _: Int ->
            Column(
                rs.getInt("ordinal_position"),
                rs.getString("column_name"),
                rs.getString("data_type"),
            )
        }, schema, table)
        return result
    }

    fun retrieve(
        schema: String,
        table: String,
        option: ConditionBuilder.Option,
        filter: ConditionBuilder.Filter
    ): List<Map<String, Any>> {
        val columns = retrieveColumns(schema, table)
        val q = "select ${columns.joinToString(", ") { it.columnName }} from $schema.$table"
        val conditions = mutableListOf<String>()
        return listOf()
    }

    fun create(schema: String, table: String, data: MutableMap<String, Any>) {
        val columns = retrieveColumns(schema, table)
        data["id"] = IdGenerator.snowflakeId()
        val state = mapOf("created_at" to Date().time, "uuid" to UUID.randomUUID().toString())
        val objectMapper = ObjectMapper()
        val s = objectMapper.writeValueAsString(state)
        data["state"] = s
        val q = """
        insert into $schema.$table (${columns.joinToString(", ") { it.columnName }})
        values (${columns.joinToString(", ") { "?" }})
        """.trimIndent()
        jdbcTemplate.update(q, *columns.map { data[it.columnName] }.toTypedArray())
    }

    fun retrieveByID(schema: String, table: String, id: Long, uuid: String): MutableMap<String, Any> {
        val columns = retrieveColumns(schema, table)
        val q = """
        select ${columns.joinToString(", ") { it.columnName }}
        from $schema.$table
        where id = ? and state->>'uuid' = '%s' limit 1
        """.trimIndent()
        val result = jdbcTemplate.queryForObject(q, { rs: ResultSet, _: Int ->
            val map = mutableMapOf<String, Any>()
            for (column in columns) {
                map[column.columnName] = rs.getObject(column.columnName)
            }
            map
        }, id, uuid)
        return result!!
    }

    fun update(schema: String, table: String, data: MutableMap<String, Any>) {
        val columns = retrieveColumns(schema, table)
        val conditions = mutableListOf<String>()
        conditions.add("state = state || jsonb_build_object('updated_at', ?)")
        for (column in columns) {
            if (data.containsKey(column.columnName)) {
                conditions.add("${column.columnName} = ?")
            }
        }
        val q = """
        update $schema.$table
        set ${conditions.joinToString(", ")}
        where id = ?
        """.trimIndent()
        jdbcTemplate.update(q, Date().time, *columns.map { data[it.columnName] }.toTypedArray(), data["id"])
    }

}
