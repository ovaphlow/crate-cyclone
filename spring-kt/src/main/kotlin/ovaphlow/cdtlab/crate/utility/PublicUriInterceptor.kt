package ovaphlow.cdtlab.crate.utility

import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Component
import org.springframework.web.servlet.HandlerInterceptor
import java.util.*

@Component
class PublicUriInterceptor : HandlerInterceptor {

    var logger: Logger = LoggerFactory.getLogger(PublicUriInterceptor::class.java)

    override fun preHandle(request: HttpServletRequest, response: HttpServletResponse, handler: Any): Boolean {
        logger.info("${request.method} ${request.requestURI}")
        return if (Constants.PUBLIC_URIS.contains(request.requestURI)) {
            true
        } else {
            var token = request.getHeader("Authorization")
            if (Objects.isNull(token)) {
                response.status = 401
                return false
            }
            token = token.replace("Bearer ", "")
            logger.info("token $token")
            return true
        }
    }

}
