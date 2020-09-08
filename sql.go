package jda

import (
	"reflect"
	"database/sql"
)

func SelectOneFromSqlTable(
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
		err = rows.Scan(append(idInter, outputFields)...)
		if err != nil {
			l.Error(err.Error())
			l.Error("Error in scan fields to SQL select query row")
			return true, l.ErrorQueue
		}
		return true, nil
	}

	return false, nil
}

func SelectFromSqlTable(
	database *sql.DB,
	templateInterface interface{},
	tableName string,
	queryExpression string,
	args ...interface{},
) ([]int, [][]interface{}, error) {
	l := GetLogger()

	if tableName == "" {
		l.Error("Table name is null")
		return false, l.ErrorQueue
	}

	expression := "SELECT * FROM "+tableName
	if queryExpression != "" {
		expression = expression+" WHERE "+queryExpression
	}

	rows, err := database.Query(expression, args...)
	if err != nil {
		l.Error(err.Error())
		l.Error("Error in query entity from SQL database")
		return nil, nil, l.ErrorQueue
	}

	structValue := reflect.Indirect(reflect.ValueOf(templateInterface))

	length := structValue.NumField()
	if length == 0 {
		l.Error("Empty interface")
		return nil, nil, l.ErrorQueue
	}

	var templateFields []interface{} = nil
	for i := 0; i < length; i = i + 1 {
		templateFields = append(
			templateFields,
			structValue.Field(i).Addr().Interface(),
		)
	}
	
	var outputFieldsArray [][]interface{} = nil
	var outputIds []int 
	for rows.Next() {
		buffer := make([]interface{}, len(templateFields))
		copy(buffer, templateFields)

		var id int = 0
		var idTemplate interface{} = &id
		err = rows.Scan(append(idTemplate, buffer)...)
		if err != nil {
			l.Error(err.Error())
			l.Error("Error in scan fields to SQL select query row")
			return nil, nil, l.ErrorQueue
		}
		
		outputIds = append(outputIds, id)
		outputFieldsArray = append(outputFieldsArray, buffer)
	}
	
	return outputIds, outputFieldsArray, nil
}

func InsertIntoSqlTable(database *sql.DB, s interface{}, tableName string) error {
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

func CreateSqlTable(database *sql.DB, s interface{}, tableName string) error {
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
