package ovaphlow.cdtlab.crate.table

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.CrossOrigin
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestMethod
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.bind.annotation.RestController

@RestController
@RequestMapping(path = ["/crate-api/table"])
@CrossOrigin
class TableController {

    @Autowired
    lateinit var tableRepository: TableRepository

    @RequestMapping(path = [""], method = [RequestMethod.GET])
    fun get(@RequestParam(value = "schema", defaultValue = "") schema: String): ResponseEntity<Any> {
        if (schema.isEmpty()) {
            return ResponseEntity.badRequest().body(mapOf("error" to "schema is required"))
        }
        val result = tableRepository.retrieve(schema)
        return ResponseEntity.ok().body(result)
    }

}
