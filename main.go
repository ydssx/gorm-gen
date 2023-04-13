package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

func main() {
	sql := `CREATE TABLE users (
		id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
		username VARCHAR(50) NOT NULL COMMENT '用户名',
		password VARCHAR(100) NOT NULL COMMENT '密码',
		email VARCHAR(50) NOT NULL COMMENT '邮箱',
		phone VARCHAR(20) NOT NULL COMMENT '电话',
		PRIMARY KEY (id),
		UNIQUE KEY idx_users_username (username),
		UNIQUE KEY idx_users_email (email),
		UNIQUE KEY idx_users_phone (phone)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='用户表';`
	table, err := ParseSQL(sql)
	if err != nil {
		log.Print("failed to parse sql:", err)
		return
	}

	modelTemplate, _ := ioutil.ReadFile("model.tmpl")
	funcMap := template.FuncMap{"Title": strings.Title, "Lower": toLowerFirst}
	// 解析模板
	tmpl, err := template.New("model").Funcs(funcMap).Parse(string(modelTemplate))
	if err != nil {
		fmt.Println("failed to parse template:", err)
		return
	}
	// 将模型转换为模板需要的数据
	data := map[string]interface{}{
		"TableName":    table.Name,
		"TableComment": table.Comment,
		"Fields":       table.Fields,
		"Name":         GetSingularTableName(table.Name),
	}
	// 将模板应用到数据上，生成代码
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Println("failed to generate code:", err)
		return
	}

	// 格式化生成的代码
	formattedCode, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Println("failed to format code:", err)
		return
	}

	// 将生成的代码写入文件
	if err := ioutil.WriteFile(fmt.Sprintf("%s.go", strings.ToLower(table.Name)), formattedCode, 0644); err != nil {
		fmt.Println("failed to write code to file:", err)
		return
	}

	fmt.Println("code generation succeeded!")
}
