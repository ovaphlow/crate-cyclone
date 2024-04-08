package ovaphlow.cdtlab.crate.event

import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.jdbc.core.JdbcTemplate
import org.springframework.stereotype.Repository
import ovaphlow.cdtlab.crate.utility.ConditionBuilder
import java.sql.ResultSet
import java.util.*

@Repository
class EventRepository(private val jdbcTemplate: JdbcTemplate) {

    val columns = arrayOf("id", "relation_id", "reference_id", "tags", "detail", "time")
    val logger: Logger = LoggerFactory.getLogger(EventRepository::class.java)

    fun defaultFilter(option: ConditionBuilder.Option, filter: ConditionBuilder.Filter): List<Event> {
        val columnListString = columns.joinToString(", ")
        var q = "select $columnListString from events"
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        val builder = ConditionBuilder()
        if (Objects.nonNull(filter.equal)) {
            val (c, p) = builder.equal(filter.equal)
            conditions += c
            params += p
        }
        if (Objects.nonNull(filter.objectContain)) {
            val (c, p) = builder.objectContain(filter.objectContain)
            conditions += c
            params += p
        }
        if (Objects.nonNull(filter.arrayContain)) {
            val (c, p) = builder.arrayContain(filter.objectContain)
            conditions += c
            params += p
        }
        if (Objects.nonNull(filter.like)) {
            val (c, p) = builder.like(filter.like)
            conditions += c
            params += p
        }
        if (Objects.nonNull(filter.objectLike)) {
            val (c, p) = builder.objectLike(filter.objectLike)
            conditions += c
            params += p
        }
        if (Objects.nonNull(filter.inList)) {
            val (c, p) = builder.inList(filter.inList)
            conditions += c
            params += p
        }
        if (Objects.nonNull(filter.lesser)) {
            val (c, p) = builder.lesser(filter.lesser)
            conditions += c
            params += p
        }
        if (Objects.nonNull(filter.greater)) {
            val (c, p) = builder.greater(filter.greater)
            conditions += c
            params += p
        }
        if (conditions.size > 0) {
            q += conditions.joinToString(" and ", " where ")
        }
        q += " order by id desc limit ${option.skip}, ${option.take}"
        logger.info(q)
        val result = jdbcTemplate.query(q, { rs: ResultSet, _: Int ->
            Event(
                rs.getLong("id"),
                rs.getLong("relation_id"),
                rs.getLong("reference_id"),
                rs.getString("tags"),
                rs.getString("detail"),
                rs.getDate("time")
            )
        }, *params.toTypedArray())
        return result
    }

}
