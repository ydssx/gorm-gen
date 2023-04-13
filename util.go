package main

import (
	"strings"

	"github.com/jinzhu/inflection"
)

func getType(token string) string {
	token = strings.ToLower(token)
	switch {
	case strings.HasPrefix(token, "bigint"):
		return "int64"
	case strings.HasPrefix(token, "int"):
		return "int"
	case strings.HasPrefix(token, "tinyint(1)"):
		return "bool"
	case strings.HasPrefix(token, "varchar"), strings.HasPrefix(token, "text"), strings.HasPrefix(token, "char"):
		return "string"
	case strings.HasPrefix(token, "decimal"), strings.HasPrefix(token, "double"):
		return "float64"
	case strings.HasPrefix(token, "float"):
		return "float32"
	case strings.HasPrefix(token, "timestamp"), strings.HasPrefix(token, "datetime"):
		return "time.Time"
	default:
		return token
	}
}

func GetSingularTableName(tableName string) string {
	// 1. 先将表名转换为驼峰式命名
	tableName = strings.ReplaceAll(tableName, "_", " ")
	tableName = strings.Title(tableName)
	tableName = strings.ReplaceAll(tableName, " ", "")

	// 2. 使用 inflection 库将驼峰式命名的表名转换为单数形式
	singularName := inflection.Singular(tableName)

	return singularName
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[0:1]) + s[1:]
}
