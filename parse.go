package main

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

type Table struct {
	Name    string
	Comment string
	Fields  []Field
}

type Field struct {
	Name     string
	Type     string
	Primary  bool
	Unique   bool
	Nullable bool
	Default  interface{}
	Comment  string
	Tag      string
}

var skipFields = []string{"id", "created_at", "updated_at", "deleted_at"}

func ParseSQL(sql string) (*Table, error) {
	table := &Table{}
	fields := []Field{}
	// Extract table name
	tableNameStart := strings.Index(sql, "CREATE TABLE ") + 13
	tableNameEnd := strings.Index(sql[tableNameStart:], " ")
	table.Name = sql[tableNameStart : tableNameStart+tableNameEnd]
	table.Name = strings.ReplaceAll(table.Name, "`", "")

	// Extract table comment
	if strings.Contains(sql, "COMMENT='") {
		commentStart := strings.Index(sql, "COMMENT='") + 9
		commentEnd := strings.Index(sql[commentStart:], "'")
		table.Comment = sql[commentStart : commentStart+commentEnd]
	}

	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CREATE TABLE") || strings.HasPrefix(line, "KEY") || strings.HasPrefix(line, ")") || strings.HasPrefix(line, "CONSTRAINT") {
			continue
		} else if strings.HasPrefix(line, ") ENGINE=") {
			break
		} else if strings.HasPrefix(line, "PRIMARY KEY") {
			pk := getPrimaryKey(line)
			for i, f := range fields {
				if SliceContain(pk, f.Name) {
					fields[i].Primary = true
				}
			}
		} else if strings.HasPrefix(line, "UNIQUE KEY") {
			idx := getIndex(line)
			for i, f := range fields {
				if f.Name == idx {
					fields[i].Unique = true
				}
			}
		} else {
			field := getField(line)
			// if SliceContain(skipFields, field.Name) {
			// 	continue
			// }
			// field.Tag = generateStructTag(field)
			fields = append(fields, field)
		}
	}
	for i, field := range fields {
		fields[i].Tag = generateStructTag(field)
	}
	table.Fields = fields
	return table, nil
}

func getTableName(line string) string {
	tokens := strings.Split(line, " ")
	return strings.TrimSuffix(tokens[2], " (")
}

func getTableComment(line string) string {
	comment := ""
	if strings.Contains(line, "COMMENT") {
		start := strings.Index(line, "COMMENT '") + 9
		end := strings.Index(line[start:], "'") + start
		comment = line[start:end]
	}
	return comment
}

func getPrimaryKey(line string) []string {
	start := strings.Index(line, "(") + 1
	end := strings.Index(line, ")")
	x := line[start:end]
	x = strings.ReplaceAll(x, "`", "")
	return strings.Split(x, ",")
}

func getIndex(line string) string {
	start := strings.Index(line, "(") + 1
	end := strings.Index(line, ")")
	return line[start:end]
}

func getField(line string) Field {
	field := Field{}
	tokens := strings.Split(line, " ")
	field.Name = strings.TrimSuffix(tokens[0], ",")
	field.Name = strings.ReplaceAll(field.Name, "`", "")
	field.Type = getType(tokens[1])
	if strings.Contains(line, "NOT NULL") {
		field.Nullable = false
	} else {
		field.Nullable = true
	}
	if strings.Contains(line, "DEFAULT") {
		start := strings.Index(line, "DEFAULT ") + 8
		fval := strings.TrimRight(strings.Split(line[start:], " ")[0], ",")
		field.Default = pareDefaultValue(field.Type, fval)
	}
	if strings.Contains(line, "COMMENT") {
		start := strings.Index(line, "COMMENT '") + 9
		end := strings.Index(line[start:], "'") + start
		field.Comment = line[start:end]
	}
	return field
}

// 生成模型tag
func generateStructTag(field Field) (tag string) {
	// fieldStr := fmt.Sprintf("%s %s", field.Name, field.Type)
	tags := []string{fmt.Sprintf("column:%s", field.Name)}
	if field.Primary {
		tags = append(tags, "primaryKey")
	}
	if field.Unique {
		tags = append(tags, "unique")
	}
	if !field.Nullable {
		tags = append(tags, "not null")
	}
	if field.Default != "" && field.Default != nil {
		tags = append(tags, fmt.Sprintf("default:%v", field.Default))
	}
	tag += fmt.Sprintf("`json:\"%s\" gorm:\"%s\"`", field.Name, strings.Join(tags, ";"))
	return tag
}

func ParseSQL1(sql string) (*Table, error) {
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		panic(err)
	}

	createStmt, ok := stmt.(*sqlparser.DDL)
	if !ok || createStmt.Action != "create" {
		return nil, fmt.Errorf("invalid create statement")
	}
	table := new(Table)
	table.Name = createStmt.NewName.Name.String()
	// fmt.Printf("Table Name: %s\n", tableName)
	columns := createStmt.TableSpec.Columns
	for _, col := range columns {
		field := Field{}
		field.Comment = string(col.Type.Comment.Val)
		if col.Type.Default != nil {
			field.Default = string(col.Type.Default.Val)
		}
		field.Name = col.Name.String()
		field.Nullable = bool(col.Type.NotNull)
		field.Type = col.Type.Type
		table.Fields = append(table.Fields, field)
	}
	return table, nil
}

func findPrimaryKey(table Table) string {
	for _, v := range table.Fields {
		if v.Primary {
			return v.Name
		}
	}
	return ""
}
