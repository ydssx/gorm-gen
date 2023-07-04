package main

import (
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/jinzhu/inflection"
)

func getType(token string) string {
	token = strings.ToLower(token)
	switch {
	case strings.HasPrefix(token, "bigint"):
		return "int64"
	case strings.HasPrefix(token, "int"), strings.HasPrefix(token, "tinyint"), strings.HasPrefix(token, "smallint"):
		return "int"
	case strings.HasPrefix(token, "tinyint(1)"):
		return "bool"
	case strings.HasPrefix(token, "varchar"), strings.HasPrefix(token, "text"), strings.HasPrefix(token, "char"), strings.HasPrefix(token, "longtext"):
		return "string"
	case strings.HasPrefix(token, "decimal"), strings.HasPrefix(token, "double"):
		return "float64"
	case strings.HasPrefix(token, "float"):
		return "float32"
	case strings.HasPrefix(token, "timestamp"), strings.HasPrefix(token, "datetime"):
		return "time.Time"
	case strings.HasPrefix(token, "json"):
		return "json.RawMessage"
	default:
		return token
	}
}

func pareDefaultValue(ftype, fval string) (v interface{}) {
	if strings.ToLower(fval) == "null" {
		return fval
	}
	switch ftype {
	case "int64", "int", "int32":
		v, _ = strconv.ParseInt(fval, 10, 64)
	case "float64", "float32":
		vf, _ := strconv.ParseFloat(fval, 64)
		v = math.Round(vf*100) / 100
	default:
		return fval
	}
	return
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

// 下划线转驼峰
func UnderscoreToCamelCase(s string) string {
	var (
		b  strings.Builder
		up bool
	)

	for _, c := range s {
		if c == '_' {
			up = true
			continue
		}

		if up {
			b.WriteRune(unicode.ToUpper(c))
			up = false
		} else {
			b.WriteRune(c)
		}
	}

	return b.String()
}

func SliceContain(s []string, elem string) bool {
	for _, v := range s {
		if v == elem {
			return true
		}
	}
	return false
}

func DirExists(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	}
	return false
}
