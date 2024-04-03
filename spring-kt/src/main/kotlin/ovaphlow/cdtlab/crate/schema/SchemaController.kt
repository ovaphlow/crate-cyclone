package ovaphlow.cdtlab.crate.schema

import org.springframework.beans.factory.annotation.Autowired
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.CrossOrigin
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestMethod
import org.springframework.web.bind.annotation.RestController
import ovaphlow.cdtlab.crate.utility.Constants

@RestController
@RequestMapping(path = ["/crate-api/schema"])
@CrossOrigin
class SchemaController {

    @Autowired
    lateinit var schemaRepository: SchemaRepository

    @RequestMapping(path = [""], method = [RequestMethod.GET])
    fun get(): ResponseEntity<Any> {
        val result = schemaRepository.retrieve()
        return ResponseEntity.ok().header(Constants.HEADER_API_VERSION, "2024-01-06").body(result)
    }

}
