package ovaphlow.cdtlab.crate.subscriber

import java.util.*

open class Subscriber(
    var id: Long,
    var email: String,
    var name: String,
    var phone: String,
    var tags: String,
    var detail: String,
    var time: Date
) {

    override fun toString(): String {
        return "Subscriber(id=$id, email='$email', name='$name', phone='$phone', tags='$tags', detail='$detail', time=$time)"
    }

}
