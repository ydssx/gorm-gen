package main

import (
	"fmt"
	"strings"
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
	Default  string
	Comment  string
	Tag      string
}

func ParseSQL(sql string) (*Table, error) {
	table := &Table{}
	fields := []Field{}
	// Extract table name
	tableNameStart := strings.Index(sql, "CREATE TABLE ") + 13
	tableNameEnd := strings.Index(sql[tableNameStart:], " ")
	table.Name = sql[tableNameStart : tableNameStart+tableNameEnd]

	// Extract table comment
	commentStart := strings.Index(sql, "COMMENT='") + 9
	commentEnd := strings.Index(sql[commentStart:], "'")
	table.Comment = sql[commentStart : commentStart+commentEnd]

	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CREATE TABLE") {
			continue
		} else if strings.HasPrefix(line, ") ENGINE=") {
			break
		} else if strings.HasPrefix(line, "PRIMARY KEY") {
			pk := getPrimaryKey(line)
			for i, f := range fields {
				if f.Name == pk {
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
			field.Tag = generateStructTag(field)
			fields = append(fields, field)
		}
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

func getPrimaryKey(line string) string {
	start := strings.Index(line, "(") + 1
	end := strings.Index(line, ")")
	return line[start:end]
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
	field.Type = getType(tokens[1])
	if strings.Contains(line, "NOT NULL") {
		field.Nullable = false
	} else {
		field.Nullable = true
	}
	if strings.Contains(line, "DEFAULT") {
		start := strings.Index(line, "DEFAULT '") + 9
		end := strings.Index(line[start:], "'") + start
		field.Default = line[start:end]
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
	if field.Primary {
		tag += " `gorm:\"primaryKey\"`"
	} else {
		tags := []string{fmt.Sprintf("column:%s", field.Name)}
		if field.Unique {
			tags = append(tags, "unique")
		}
		if !field.Nullable {
			tags = append(tags, "not null")
		}
		if field.Default != "" {
			tags = append(tags, fmt.Sprintf("default:%s", field.Default))
		}
		tag += fmt.Sprintf("`json:\"%s\" gorm:\"%s\"`", field.Name, strings.Join(tags, ";"))
	}
	return tag
}
