package jda

import (
	"reflect"
	"database/sql"
)

func SelectOneFromSqlTable(
	database *sql.DB,
	tableName string,
	outputInterface *interface{},
	queryExpression string,
	args ...interface{},
) error {
	l := GetLogger()

	if tableName == "" {
		l.Error("Table name is null")
		return l.ErrorQueue
	}

	expression := "SELECT * FROM "+tableName

	if queryExpression != "" {
		expression = expression+" WHERE "+queryExpression
	}

	//...

	return nil
}

func InsertIntoSqlTable(database *sql.DB, tableName string, s interface{}) error {
	l := GetLogger()

	if tableName == "" {
		l.Error("Table name is null")
		return l.ErrorQueue
	}

	v := reflect.ValueOf(s)

	length := v.NumField()
	if length == 0 {
		l.Error("Empty interface")
		return l.ErrorQueue
	}

	var fieldsToInsert []interface{} = nil

	expression := "INSERT INTO "+tableName+" ("

	for i := 0; i < length; i = i + 1 {
		fieldsToInsert = append(fieldsToInsert, v.Field(i).Interface())
		fieldName := v.Type().Field(i).Name

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

func CreateSqlTable(database *sql.DB, tableName string, s interface{}) error {
	l := GetLogger()

	if tableName == "" {
		l.Error("Table name is null")
		return l.ErrorQueue
	}

	v := reflect.ValueOf(s)

	length := v.NumField()
	if length == 0 {
		l.Error("Empty interface")
		return l.ErrorQueue
	}

	expression := "CREATE TABLE IF NOT EXISTS "+tableName+" ("

	for i := 0; i < length; i = i + 1 {
		fieldName := v.Type().Field(i).Name

		//get SQL tags for each field
		tags := v.Type().Field(i).Tag.Get("sql")
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