package ovaphlow.cdtlab.crate.event

import java.util.*

open class Event(
    var id: Long,
    var relationId: Long,
    var referenceId: Long,
    var tags: String,
    var detail: String,
    var time: Date
)
