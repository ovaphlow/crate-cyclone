package ovaphlow.cdtlab.crate

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class CrateApplication

fun main(args: Array<String>) {
	runApplication<CrateApplication>(*args)
}
