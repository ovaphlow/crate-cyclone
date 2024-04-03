package ovaphlow.cdtlab.crate.utility

import org.slf4j.Logger
import org.slf4j.LoggerFactory

class QueryBuilder {

    private val logger: Logger = LoggerFactory.getLogger(QueryBuilder::class.java)

    private var columns: String = ""
    private var schema: String = ""
    private var table: String = ""
    private var conditions: List<String> = mutableListOf()
    private var params: List<String> = mutableListOf()
    private var order: String = ""
    private var take: Int = 10
    private var skip: Long = 0

    fun select(columns: String): QueryBuilder {
        this.columns = columns
        return this
    }

    fun from(schema: String, table: String): QueryBuilder {
        this.schema = schema
        this.table = table
        return this
    }

    fun equal(equal: List<String>): QueryBuilder {
        if (equal.isEmpty() || equal.size % 2 != 0) {
            logger.debug("equal 参数错误")
            return this
        }
        for (i in equal.indices step 2) {
            this.conditions += "${equal[i]} = ?"
            this.params += equal[i + 1]
        }
        return this
    }

    fun notEqual(notEqual: List<String>): QueryBuilder {
        if (notEqual.isEmpty() || notEqual.size % 2 != 0) {
            logger.debug("notEqual 参数错误")
            return this
        }
        for (i in notEqual.indices step 2) {
            this.conditions += "${notEqual[i]} != ?"
            this.params += notEqual[i + 1]
        }
        return this
    }

    fun inList(inList: List<String>): QueryBuilder {
        if (inList.isEmpty() || inList.size % 2 != 0) {
            logger.debug("inList 参数错误")
            return this
        }
        val c: List<String> = List(inList.size - 1) { "?" }
        this.conditions += "${inList[0]} in (${c.joinToString(", ")})"
        this.params += inList.subList(1, inList.size)
        return this
    }

    fun notInList(notInList: List<String>): QueryBuilder {
        if (notInList.isEmpty() || notInList.size % 2 != 0) {
            logger.debug("notInList 参数错误")
            return this
        }
        val c: List<String> = List(notInList.size - 1) { "?" }
        this.conditions += "${notInList[0]} not in (${c.joinToString(", ")})"
        this.params += notInList.subList(1, notInList.size)
        return this
    }

    fun objectContain(objectContain: List<String>): QueryBuilder {
        if (objectContain.isEmpty() || objectContain.size % 2 != 0) {
            logger.debug("objectContain 参数错误")
            return this
        }
        for (i in objectContain.indices step 3) {
            this.conditions += "${objectContain[i]} @> jsonb_build_object('${objectContain[i + 1]}', ?)"
            this.params += objectContain[i + 2]
        }
        return this
    }

    fun arrayContain(arrayContain: List<String>): QueryBuilder {
        if (arrayContain.isEmpty() || arrayContain.size % 2 != 0) {
            logger.debug("arrayContain 参数错误")
            return this
        }
        for (i in arrayContain.indices step 2) {
            this.conditions += "${arrayContain[i]} @> jsonb_build_array(?)"
            this.params += arrayContain[i + 1]
        }
        return this
    }

    fun like(like: List<String>): QueryBuilder {
        if (like.isEmpty() || like.size % 2 != 0) {
            logger.debug("like 参数错误")
            return this
        }
        for (i in like.indices step 2) {
            this.conditions += "position(? in ${like[i]}) > 0"
            this.params += like[i + 1]
        }
        return this
    }

    fun greater(greater: List<String>): QueryBuilder {
        if (greater.isEmpty() || greater.size % 2 != 0) {
            logger.debug("greater 参数错误")
            return this
        }
        for (i in greater.indices step 2) {
            this.conditions += "${greater[i]} >= ?"
            this.params += greater[i + 1]
        }
        return this
    }

    fun lesser(lesser: List<String>): QueryBuilder {
        if (lesser.isEmpty() || lesser.size % 2 != 0) {
            logger.debug("lesser 参数错误")
            return this
        }
        for (i in lesser.indices step 2) {
            this.conditions += "${lesser[i]} <= ?"
            this.params += lesser[i + 1]
        }
        return this
    }

    fun orderBy(order: String): QueryBuilder {
        this.order = order
        return this
    }

    fun limit(take: Int, skip: Long): QueryBuilder {
        this.take = take
        this.skip = skip
        return this
    }

}
