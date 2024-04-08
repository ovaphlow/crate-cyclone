package ovaphlow.cdtlab.crate.subscriber

import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository
import java.sql.ResultSet

@Repository
class SubscriberRepository(private val jdbcTemplate: JdbcTemplate) {

    val columns = arrayOf("id", "email", "name", "phone", "tags", "detail", "time")
    val logger: Logger = LoggerFactory.getLogger(SubscriberRepository::class.java)

    fun retrieveByUsername(username: String): List<Subscriber> {
        val columnListString = columns.joinToString(", ")
        val q = """
        select $columnListString from subscribers
        where email = ? or name = ? or phone = ?
        """.trimIndent()
        val result = jdbcTemplate.query(q, { rs: ResultSet, _: Int ->
            Subscriber(
                rs.getLong("id"),
                rs.getString("email"),
                rs.getString("name"),
                rs.getString("phone"),
                rs.getString("tags"),
                rs.getString("detail"),
                rs.getDate("time")
            )
        }, username, username, username)
        return if (result.isEmpty()) listOf() else result
    }

}
