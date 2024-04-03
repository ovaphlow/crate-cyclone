package ovaphlow.cdtlab.crate.event

import jakarta.servlet.http.HttpServletRequest
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.CrossOrigin
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestMethod
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.bind.annotation.RestController
import ovaphlow.cdtlab.crate.utility.ConditionBuilder
import ovaphlow.cdtlab.crate.utility.Constants

@RestController
@RequestMapping(path = ["/crate-api/event"])
@CrossOrigin
class EventController {

    @Autowired
    lateinit var eventRepository: EventRepository

    @RequestMapping(path = [""], method = [RequestMethod.GET])
    fun filter(
        @RequestParam(name = "option", defaultValue = "") option: String,
        @RequestParam(name = "skip", defaultValue = "0") skip: Long,
        @RequestParam(name = "take", defaultValue = "10") take: Int,
        request: HttpServletRequest
    ): ResponseEntity<Any> {
        if (option == "default") {
            val options = ConditionBuilder.Option(skip, take)
            val filters = ConditionBuilder.Filter(
                request.getParameter("equal")?.split(",") ?: listOf(),
                request.getParameter("object-contain")?.split(",") ?: listOf(),
                request.getParameter("array-contain")?.split(",") ?: listOf(),
                request.getParameter("like")?.split(",") ?: listOf(),
                request.getParameter("object-like")?.split(",") ?: listOf(),
                request.getParameter("in")?.split(",") ?: listOf(),
                request.getParameter("lesser")?.split(",") ?: listOf(),
                request.getParameter("greater")?.split(",") ?: listOf(),
            )
            val result = eventRepository.defaultFilter(options, filters)
            val response: List<MutableMap<String, Any?>> = result.map {
                mutableMapOf(
                    "id" to it.id,
                    "relationId" to it.relationId,
                    "referenceId" to it.referenceId,
                    "tags" to it.tags,
                    "detail" to it.detail,
                    "time" to it.time,
                    "_id" to it.id.toString(),
                    "_relationId" to it.relationId.toString(),
                    "_referenceId" to it.referenceId.toString(),
                )
            }
            return ResponseEntity.ok().header(Constants.HEADER_API_VERSION, "2024-01-06").body(response)
        }
        return ResponseEntity.status(406).body(mapOf("message" to "invalid option"))
    }
}
