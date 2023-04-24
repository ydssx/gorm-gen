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
	sql, _ := ioutil.ReadFile("sql.text")
	table, err := ParseSQL(string(sql))
	if err != nil {
		log.Print("failed to parse sql:", err)
		return
	}

	modelTemplate, _ := ioutil.ReadFile("model.tmpl")
	funcMap := template.FuncMap{"Title": strings.Title, "Lower": toLowerFirst, "CamelCase": UnderscoreToCamelCase}
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
	if err := ioutil.WriteFile(fmt.Sprintf("model/%s.go", strings.ToLower(table.Name)), formattedCode, 0644); err != nil {
		fmt.Println("failed to write code to file:", err)
		return
	}

	fmt.Println("code generation succeeded!")
}
