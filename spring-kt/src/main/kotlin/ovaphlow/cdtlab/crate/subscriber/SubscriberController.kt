package ovaphlow.cdtlab.crate.subscriber

import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import com.fasterxml.jackson.module.kotlin.readValue
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.beans.factory.annotation.Value
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.*
import ovaphlow.cdtlab.crate.utility.Constants
import java.security.MessageDigest

@RestController
@RequestMapping(path = ["/crate-api/subscriber"])
@CrossOrigin
class SubscriberController {

    @Value("\${CRATE_JWT_KEY}")
    lateinit var jwtKey: String

    @Autowired
    lateinit var subscriberRepository: SubscriberRepository

    class SignInRequest(
        val username: String?,
        val password: String?
    )

    class SignUpRequest(
        val email: String?,
        val name: String?,
        val phone: String?,
        val password: String?,
    )

    @RequestMapping(path = ["/sign-in"], method = [RequestMethod.POST])
    fun signIn(@RequestBody body: SignInRequest): ResponseEntity<Any> {
        println(body)
        if (body.username.isNullOrEmpty() || body.password.isNullOrEmpty()) {
            return ResponseEntity.badRequest()
                .header(Constants.HEADER_API_VERSION, "2024-01-06")
                .body(mapOf("message" to "参数错误"))
        }
        val subscribers = subscriberRepository.retrieveByUsername(body.username)
        if (subscribers.isEmpty()) {
            return ResponseEntity.status(404)
                .header(Constants.HEADER_API_VERSION, "2024-01-06")
                .body(mapOf("message" to "用户名或密码错误"))
        }
        if (subscribers.size > 1) {
            return ResponseEntity.status(500)
                .header(Constants.HEADER_API_VERSION, "2024-01-06")
                .body(mapOf("message" to "multiple subscribers found"))
        }
        subscribers.forEach(::println)
        val detail: Map<String, Any> = jacksonObjectMapper().readValue(subscribers[0].detail)
        val messageDigest = MessageDigest.getInstance("SHA-256")
        messageDigest.update(detail["salt"].toString().toByteArray())
        val bytes = messageDigest.digest(body.password.toByteArray())
        val hex = bytes.joinToString("") { "%02x".format(it) }
        println(hex)
        println(detail["sha"])
        return ResponseEntity.ok()
            .header(Constants.HEADER_API_VERSION, "2024-01-06")
            .body(mapOf("message" to "sign-in"))
    }

}
