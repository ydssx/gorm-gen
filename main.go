package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DataBase struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	Output string   `yaml:"output"`
	Tables []string `yaml:"tables"`
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.yaml", "path to config file")
	flag.Parse()

	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(config.Output)
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path: %s", err))
	}

	// 如果目录不存在，创建目录
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		if err := os.MkdirAll(absPath, 0755); err != nil {
			panic(fmt.Errorf("failed to create directory: %s", err))
		}
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DataBase.Username, config.DataBase.Password, config.DataBase.Host, config.DataBase.Port, config.DataBase.Name)
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		// 处理连接错误
		log.Fatalf("failed to connect database: %v", err)
	}

	for _, tableName := range config.Tables {
		var createSQL string
		if err := db.Raw("SHOW CREATE TABLE "+tableName).Row().Scan(&tableName, &createSQL); err != nil {
			// 处理错误
			log.Fatalf("failed to get sql: %v", err)
		}

		generate(createSQL, absPath)
	}
}

func generate(createSQL, outPath string) {
	table, err := ParseSQL(createSQL)
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
	log.Print(buf.String())
	// 将生成的代码写入文件
	filename := filepath.Join(outPath, strings.ToLower(table.Name)+".go")
	log.Print("filepath:" + filename)
	log.Print("code:" + string(formattedCode))
	if err := ioutil.WriteFile(filename, formattedCode, 0644); err != nil {
		fmt.Println("failed to write code to file:", err)
		return
	}
	s := color.BlueString("[table %s]", table.Name)
	fmt.Printf("%s:code generation succeeded!\n", s)
}
