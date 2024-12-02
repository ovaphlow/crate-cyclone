package dbutil

type Repo interface {
	// Create inserts a new record into the specified table.
	//
	// Parameters:
	// - st: schema and table, formatted as "schema.table"
	// - d: data to be inserted
	//
	// Returns:
	// - error: error information
	Create(st string, d map[string]interface{}) error

	// Get retrieves records from the specified table based on conditions.
	//
	// Parameters:
	// - st: schema and table, formatted as "schema.table"
	// - c: columns to retrieve, e.g., ["id", "name"]
	// - f: filter conditions, e.g., [["equal", "name", "John Doe"], ["in", "id", "1a", "1b"]]
	// - l: additional clauses, e.g., "order by id desc limit 20 offset 0"
	//
	// Returns:
	// - []map[string]interface{}: retrieved records
	// - error: error information
	Get(st string, c []string, f [][]string, l string) ([]map[string]interface{}, error)

	// Update modifies records in the specified table based on conditions.
	//
	// Parameters:
	// - st: schema and table, formatted as "schema.table"
	// - d: data to be updated
	// - w: WHERE condition, e.g., "id='1a'"
	//
	// Returns:
	// - error: error information
	Update(st string, d map[string]interface{}, w string) error

	// Remove deletes records from the specified table based on conditions.
	//
	// Parameters:
	// - st: schema and table, formatted as "schema.table"
	// - w: WHERE condition, e.g., "id='1a'"
	//
	// Returns:
	// - error: error information
	Remove(st string, w string) error
}
