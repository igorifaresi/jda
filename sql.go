package jda

import (
	"reflect"
	"database/sql"
)

func SqlSelectOne(
	database *sql.DB,
	id *int,
	outputInterface interface{},
	tableName string,
	queryExpression string,
	args ...interface{},
) (bool, error) {
	l := GetLogger()

	if tableName == "" {
		l.Error("Table name is null")
		return false, l.ErrorQueue
	}

	expression := "SELECT * FROM "+tableName
	if queryExpression != "" {
		expression = expression+" WHERE "+queryExpression
	}
	l.Log(expression)

	rows, err := database.Query(expression, args...)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in query entity from SQL database")
		return false, l.ErrorQueue
	}

	structValue := reflect.Indirect(reflect.ValueOf(outputInterface))

	length := structValue.NumField()
	if length == 0 {
		l.Error("Empty interface")
		return false, l.ErrorQueue
	}
	
	var outputFields []interface{} = nil
	for i := 0; i < length; i = i + 1 {
		outputFields = append(
			outputFields,
			structValue.Field(i).Addr().Interface(),
		)
	}

	*id = 0
	var idInter interface{} = id
	if rows.Next() {
		err = rows.Scan(append(outputFields, idInter)...)
		if err != nil {
			l.Error(err.Error())
			l.Error("Error in scan fields to SQL select query row")
			return true, l.ErrorQueue
		}
		return true, nil
	}

	return false, nil
}

func SqlSelect(
	database *sql.DB,
	templateInterface interface{},
	tableName string,
	queryExpression string,
	args ...interface{},
) ([]int, []interface{}, error) {
	l := GetLogger()

	if tableName == "" {
		l.Error("Table name is null")
		return nil, nil, l.ErrorQueue
	}

	expression := "SELECT * FROM "+tableName
	if queryExpression != "" {
		expression = expression+" WHERE "+queryExpression
	}
	l.Log(expression)

	rows, err := database.Query(expression, args...)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in query entity from SQL database")
		return nil, nil, l.ErrorQueue
	}

	var outputFieldsArray []interface{} = nil
	var outputIds []int
	for rows.Next() {
		rowInter := templateInterface
		rowValue := reflect.Indirect(reflect.ValueOf(rowInter))

		var rowFields []interface{} = nil
		length := rowValue.NumField()
		for i := 0; i < length; i = i + 1 {
			rowFields = append(
				rowFields,
				rowValue.Field(i).Addr().Interface(),
			)
		}

		var id int = 0
		var idTemplate interface{} = &id
		err = rows.Scan(append([]interface{}{idTemplate}, rowFields...)...)
		if err != nil {
			l.Error(err.Error())
			l.Error("Error in scan fields to SQL select query row")
			return nil, nil, l.ErrorQueue
		}

		outputIds = append(outputIds, id)
		outputFieldsArray = append(
			outputFieldsArray,
			reflect.Indirect(reflect.ValueOf(rowInter)).Interface(),
		)
	}

	return outputIds, outputFieldsArray, nil
}

func SqlInsert(database *sql.DB, s interface{}, tableName string) error {
	l := GetLogger()

	if tableName == "" {
		l.Error("Table name is null")
		return l.ErrorQueue
	}

	structValue := reflect.ValueOf(s)

	length := structValue.NumField()
	if length == 0 {
		l.Error("Empty interface")
		return l.ErrorQueue
	}

	var fieldsToInsert []interface{} = nil

	expression := "INSERT INTO "+tableName+" ("

	for i := 0; i < length; i = i + 1 {
		fieldsToInsert = append(fieldsToInsert, structValue.Field(i).Interface())
		fieldName := structValue.Type().Field(i).Name

		expression = expression+fieldName
		if i != (length - 1) {
			expression = expression+", "
		}
	}

	expression = expression+") VALUES ("
	for i := 0; i < length; i = i + 1 {
		expression = expression+"?"
		if i != (length - 1) {
			expression = expression+", "
		}
	}
	expression = expression+")"
	l.Log(expression)

	statement, err := database.Prepare(expression)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in prepare SQL insert statement")
		return l.ErrorQueue
	}
	_, err = statement.Exec(fieldsToInsert...)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in execute SQL insert statement")
		return l.ErrorQueue
	}

	return nil
}

func SqlCreateTable(database *sql.DB, s interface{}, tableName string) error {
	l := GetLogger()

	if tableName == "" {
		l.Error("Table name is null")
		return l.ErrorQueue
	}

	structValue := reflect.ValueOf(s)

	length := structValue.NumField()
	if length == 0 {
		l.Error("Empty interface")
		return l.ErrorQueue
	}

	expression := "CREATE TABLE IF NOT EXISTS "+tableName+" (id INTEGER PRIMARY KEY AUTO_INCREMENT, "

	for i := 0; i < length; i = i + 1 {
		fieldName := structValue.Type().Field(i).Name

		//get SQL tags for each field
		tags := structValue.Type().Field(i).Tag.Get("sql")
		if tags == "" || tags == "-" {
			continue
		}

		expression = expression+fieldName+" "+tags
		if i != (length - 1) {
			expression = expression+", "
		}
	}
	expression = expression+")"
	l.Log(expression)

	statement, err := database.Prepare(expression)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in prepare SQL create table statement")
		return l.ErrorQueue
	}
	_, err = statement.Exec()
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in execute SQL create table statement")
		return l.ErrorQueue
	}

	return nil
}
