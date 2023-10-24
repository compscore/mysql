package check_template

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
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

	connectInterface, ok := options["connect"]
	if ok {
		connect, ok := connectInterface.(bool)
		if ok {
			o.Connect = connect
		}
	}

	tableExistsInterface, ok := options["table_exists"]
	if ok {
		tableExists, ok := tableExistsInterface.(bool)
		if ok {
			o.TableExists = tableExists
		}
	}

	rowExistsInterface, ok := options["row_exists"]
	if ok {
		rowExists, ok := rowExistsInterface.(bool)
		if ok {
			o.RowExists = rowExists
		}
	}

	regexMatchInterface, ok := options["regex_match"]
	if ok {
		regexMatch, ok := regexMatchInterface.(bool)
		if ok {
			o.RegexMatch = regexMatch
		}
	}

	substringMatchInterface, ok := options["substring_match"]
	if ok {
		substringMatch, ok := substringMatchInterface.(bool)
		if ok {
			o.SubstringMatch = substringMatch
		}
	}

	matchInterface, ok := options["match"]
	if ok {
		match, ok := matchInterface.(bool)
		if ok {
			o.Match = match
		}
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

	conn.SetConnMaxIdleTime(-1)
	conn.SetMaxOpenConns(1)

	// Check if connection can be made
	if optionsStruct.Connect {
		err = conn.Ping()
		if err != nil {
			return false, err.Error()
		}
	}

	// Check if table exists
	if optionsStruct.TableExists {
		query, err := conn.Prepare("SHOW TABLES LIKE ?")
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.Query(optionsStruct.Table)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("table does not exist: \"%s\"", optionsStruct.Table)
		}
	}

	// Check if row exists
	if optionsStruct.RowExists {
		query, err := conn.Prepare("SELECT * FROM ? LIMIT 1")
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.Query(optionsStruct.Table)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("table is empty: \"%s\"", optionsStruct.Table)
		}
	}

	// Check if field matches regex
	if optionsStruct.RegexMatch {
		regexp, err := regexp.Compile(expectedOutput)
		if err != nil {
			return false, err.Error()
		}

		query, err := conn.Prepare("SELECT ? FROM ? LIMIT 1")
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.Query(optionsStruct.Field, optionsStruct.Table)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("table is empty: \"%s\"", optionsStruct.Table)
		}

		var field string
		err = rows.Scan(&field)
		if err != nil {
			return false, err.Error()
		}

		if !regexp.MatchString(field) {
			return false, fmt.Sprintf("field does not match regex: \"%s\"", expectedOutput)
		}
	}

	// Check if field matches substring or contains substring
	if optionsStruct.SubstringMatch || optionsStruct.Match {
		query, err := conn.Prepare("SELECT ? FROM ? LIMIT 1")
		if err != nil {
			return false, err.Error()
		}
		defer query.Close()

		rows, err := query.Query(optionsStruct.Field, optionsStruct.Table)
		if err != nil {
			return false, err.Error()
		}

		if !rows.Next() {
			return false, fmt.Sprintf("table is empty: \"%s\"", optionsStruct.Table)
		}

		var field string
		err = rows.Scan(&field)
		if err != nil {
			return false, err.Error()
		}

		if optionsStruct.Match && field != expectedOutput {
			return false, fmt.Sprintf("field does not match string: \"%s\"", expectedOutput)
		}

		if optionsStruct.SubstringMatch && !strings.Contains(field, expectedOutput) {
			return false, fmt.Sprintf("field does not contain substring: \"%s\"", expectedOutput)
		}

	}

	return true, ""
}
