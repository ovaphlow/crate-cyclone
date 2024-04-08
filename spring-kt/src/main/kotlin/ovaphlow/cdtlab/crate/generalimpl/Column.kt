package ovaphlow.cdtlab.crate.generalimpl

open class Column(val ordinalPosition: Int, val columnName: String, val dataType: String) {

    override fun toString(): String {
        return "Column(ordinalPosition=$ordinalPosition, columnName='$columnName', dataType='$dataType')"
    }

}
