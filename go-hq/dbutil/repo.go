package dbutil

type RDBRepo interface {
	// Create 插入一条新记录到指定的表。
	// 参数：
	// - st：schema 和表，格式如 "schema.table"
	// - d：要插入的数据
	// 返回值：
	// - error：错误信息
	Create(st string, d map[string]interface{}) error
	// Get 根据条件从指定的表中检索记录。
	// 参数：
	// - st：schema 和表，格式如 "schema.table"
	// - c：要检索的列，例如 ["id", "name"]
	// - f：过滤条件，例如 [["equal", "name", "John Doe"], ["in", "id", "1a", "1b"]]
	// - l：附加子句，例如 "order by id desc limit 20 offset 0"
	// 返回值：
	// - []map[string]interface{}：检索到的记录
	// - error：错误信息
	Get(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error)
	// Update 根据条件修改指定表中的记录。
	// 参数：
	// - st：schema 和表，格式如 "schema.table"
	// - d：要更新的数据
	// - w：WHERE 条件，例如 "id='1a'"
	// 返回值：
	// - error：错误信息
	Update(st string, d map[string]interface{}, w string) error
	// Remove 根据条件删除指定表中的记录。
	// ��数：
	// - st：schema 和表，格式如 "schema.table"
	// - w：WHERE 条件，例如 "id='1a'"
	// 返回值：
	// - error：错误信息
	Remove(st string, w string) error
}
