package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type optionsStruct struct {
	// Database is the name of the database to connect to
	Database string `compscore:"database"`

	// Table is the name of the table to check
	Table string `compscore:"table"`

	// Field is the name of the field to check
	Field string `compscore:"field"`

	// Check if a database connection can be made
	Connect bool `compscore:"connect"`

	// Check if a table exists
	TableExists bool `compscore:"table_exists"`

	// Check if any row exists in table
	RowExists bool `compscore:"row_exists"`

	// Check if a field of row matches regex in
	RegexMatch bool `compscore:"regex_match"`

	// Check if a field of row matches substring
	SubstringMatch bool `compscore:"substring_match"`

	// Check if a field of row matches exact string
	Match bool `compscore:"match"`
}

func (o *optionsStruct) Unmarshal(options map[string]interface{}) {
	databaseInterface, ok := options["database"]
	if ok {
		database, ok := databaseInterface.(string)
		if ok {
			o.Database = database
		}
	}

	tableInterface, ok := options["table"]
	if ok {
		table, ok := tableInterface.(string)
		if ok {
			o.Table = table
		}
	}

	fieldInterface, ok := options["field"]
	if ok {
		field, ok := fieldInterface.(string)
		if ok {
			o.Field = field
		}
	}

	_, ok = options["connect"]
	if ok {
		o.Connect = true
	}

	_, ok = options["table_exists"]
	if ok {
		o.TableExists = true
	}

	_, ok = options["row_exists"]
	if ok {
		o.RowExists = true
	}

	_, ok = options["regex_match"]
	if ok {
		o.RegexMatch = true
	}

	_, ok = options["substring_match"]
	if ok {
		o.SubstringMatch = true
	}

	_, ok = options["match"]
	if ok {
		o.Match = true
	}
}

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string, options map[string]interface{}) (bool, string) {
	optionsStruct := optionsStruct{}
	optionsStruct.Unmarshal(options)

	conn, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s)/%s",
			username,
			password,
			target,
			optionsStruct.Database,
		),
	)
	if err != nil {
		return false, err.Error()
	}
	defer conn.Close()

	// Check if connection can be made
	if optionsStruct.Connect {
		err = conn.Ping()
		if err != nil {
			return false, err.Error()
		}
	}

	// Check if table exists
	if optionsStruct.TableExists {
		query, err := conn.PrepareContext(ctx, "SHOW TABLES LIKE ?")
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.QueryContext(ctx, optionsStruct.Table)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("table does not exist: \"%s\"", optionsStruct.Table)
		}
	}

	// Check if row exists
	if optionsStruct.RowExists {
		queryStr := fmt.Sprintf("SELECT * FROM %s LIMIT 1", optionsStruct.Table)
		query, err := conn.PrepareContext(ctx, queryStr)
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.QueryContext(ctx)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("table is empty: \"%s\"", optionsStruct.Table)
		}
	}

	// Check if field matches regex
	if optionsStruct.RegexMatch {
		queryStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s REGEXP ? LIMIT 1", optionsStruct.Field, optionsStruct.Table, optionsStruct.Field)
		query, err := conn.PrepareContext(ctx, queryStr)
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.QueryContext(ctx, expectedOutput)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("no results exist that match regex: \"%s\" for field: \"%s\" in table: \"%s\" in database: \"%s\"", expectedOutput, optionsStruct.Field, optionsStruct.Table, optionsStruct.Database)
		}
	}

	// Check if field matches substring
	if optionsStruct.SubstringMatch {
		queryStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s LIKE ? LIMIT 1", optionsStruct.Field, optionsStruct.Table, optionsStruct.Field)
		query, err := conn.PrepareContext(ctx, queryStr)
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.QueryContext(ctx, expectedOutput)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("no results exist that match substring: \"%s\" for field: \"%s\" in table: \"%s\" in database: \"%s\"", expectedOutput, optionsStruct.Field, optionsStruct.Table, optionsStruct.Database)
		}
	}

	// Check if field matches exact string
	if optionsStruct.Match {
		queryStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? LIMIT 1", optionsStruct.Field, optionsStruct.Table, optionsStruct.Field)
		query, err := conn.PrepareContext(ctx, queryStr)
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.QueryContext(ctx, expectedOutput)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("no results exist that match string: \"%s\" for field: \"%s\" in table: \"%s\" in database: \"%s\"", expectedOutput, optionsStruct.Field, optionsStruct.Table, optionsStruct.Database)
		}
	}

	return true, ""
}
