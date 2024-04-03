package ovaphlow.cdtlab.crate.utility

import org.slf4j.Logger
import org.slf4j.LoggerFactory

class ConditionBuilder {

    private var logger: Logger = LoggerFactory.getLogger(ConditionBuilder::class.java)

    class Option(val skip: Long, val take: Int)

    class Filter(
        val equal: List<String>,
        val objectContain: List<String>,
        val arrayContain: List<String>,
        val like: List<String>,
        val objectLike: List<String>,
        val inList: List<String>,
        val lesser: List<String>,
        val greater: List<String>
    )

    fun equal(equal: List<String>): Pair<List<String>, List<String>> {
        if (equal.isEmpty() || equal.size % 2 != 0) {
            logger.debug("equal 参数错误")
            return Pair(listOf(), listOf())
        }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        for (i in equal.indices step 2) {
            conditions += "${equal[i]} = ?"
            params += equal[i + 1]
        }
        return Pair(conditions, params)
    }

    fun notEqual(notEqual: List<String>): Pair<List<String>, List<String>> {
        if (notEqual.isEmpty() || notEqual.size % 2 != 0) {
            logger.debug("notEqual 参数错误")
            return Pair(listOf(), listOf())
        }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        for (i in notEqual.indices step 2) {
            conditions += "${notEqual[i]} != ?"
            params += notEqual[i + 1]
        }
        return Pair(conditions, params)
    }

    fun like(like: List<String>): Pair<List<String>, List<String>> {
        if (like.isEmpty() || like.size % 2 != 0) {
            logger.debug("like 参数错误")
            return Pair(listOf(), listOf())
        }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        for (i in like.indices step 2) {
            conditions += "position(? in ${like[i]})"
            params += like[i + 1]
        }
        return Pair(conditions, params)
    }

    fun greater(greater: List<String>): Pair<List<String>, List<String>> {
        if (greater.isEmpty() || greater.size % 2 != 0) {
            logger.debug("greater 参数错误")
            return Pair(listOf(), listOf())
        }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        for (i in greater.indices step 2) {
            conditions += "${greater[i]} >= ?"
            params += greater[i + 1]
        }
        return Pair(conditions, params)
    }

    fun lesser(lesser: List<String>): Pair<List<String>, List<String>> {
        if (lesser.isEmpty() || lesser.size % 2 != 0) {
            logger.debug("lesser 参数错误")
            return Pair(listOf(), listOf())
        }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        for (i in lesser.indices step 2) {
            conditions += "${lesser[i]} <= ?"
            params += lesser[i + 1]
        }
        return Pair(conditions, params)
    }

    fun inList(inList: List<String>): Pair<List<String>, List<String>> {
        if (inList.isEmpty() || inList.size < 2) {
            logger.debug("inList 参数错误")
            return Pair(listOf(), listOf())
        }
        val c: List<String> = List(inList.size - 1) { "?" }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        conditions += "${inList[0]} in (${c.joinToString(", ")})"
        params += inList.subList(1, inList.size)
        return Pair(conditions, params)
    }

    fun notInList(notInList: List<String>): Pair<List<String>, List<String>> {
        if (notInList.isEmpty() || notInList.size < 2) {
            logger.debug("notInList 参数错误")
            return Pair(listOf(), listOf())
        }
        val c: List<String> = List(notInList.size - 1) { "?" }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        conditions += "${notInList[0]} not in (${c.joinToString(", ")})"
        params += notInList.subList(1, notInList.size)
        return Pair(conditions, params)
    }

    fun objectContain(objectContain: List<String>): Pair<List<String>, List<String>> {
        if (objectContain.isEmpty() || objectContain.size % 3 != 0) {
            logger.debug("objectContain 参数错误")
            return Pair(listOf(), listOf())
        }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        for (i in objectContain.indices step 3) {
            conditions += "json_contains(${objectContain[i]}, json_object('${objectContain[i + 1]}', ?))"
            params += objectContain[i + 2]
        }
        return Pair(conditions, params)
    }

    fun arrayContain(arrayContain: List<String>): Pair<List<String>, List<String>> {
        if (arrayContain.isEmpty() || arrayContain.size % 2 != 0) {
            logger.debug("arrayContain 参数错误")
            return Pair(listOf(), listOf())
        }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        for (i in arrayContain.indices step 2) {
            conditions += "json_contains(${arrayContain[i]}, json_array(?))"
            params += arrayContain[i + 1]
        }
        return Pair(conditions, params)
    }

    fun objectLike(objectLike: List<String>): Pair<List<String>, List<String>> {
        if (objectLike.isEmpty() || objectLike.size % 3 != 0) {
            logger.debug("objectLike 参数错误")
            return Pair(listOf(), listOf())
        }
        val conditions = mutableListOf<String>()
        val params = mutableListOf<String>()
        for (i in objectLike.indices step 3) {
            conditions += "position(? in ${objectLike[i]}->>'$.${objectLike[i + 1]}')"
            params += objectLike[i + 2]
        }
        return Pair(conditions, params)
    }

}
