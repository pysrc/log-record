package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

// 配置
type Config struct {
	Host *string `json:"host"` // 服务地址
	Db   *string `json:"db"`   // 数据库地址
}

// 日志
type Log struct {
	Svc  string `json:"svc"`  // 日志系统
	Info string `json:"info"` // 日志详情
}

//错误处理
func Error(err error) {
	if err != nil {
		log.Panic(err)
		panic(err)
	}
}

// 解析配置文件
func ParseConfig(name string) *Config {
	data, err := ioutil.ReadFile(name)
	Error(err)
	var info Config
	err = json.Unmarshal(data, &info)
	Error(err)
	return &info
}

// 初始化表结构
func InitTable(db *sql.DB) {
	sql_create := `
    CREATE TABLE IF NOT EXISTS LOG_RECORD(
        log_system varchar(200) COMMENT "日志系统",
        log_info TEXT COMMENT "日志详细",
        log_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT "记录日期"
    )
    `
	stmt, err := db.Prepare(sql_create)
	Error(err)
	defer stmt.Close()
	_, err = stmt.Exec()
	Error(err)
}

// 插入日志
func Insert(db *sql.DB, svc string, info string) {
	stmt, _ := db.Prepare(`insert into log_record(log_system, log_info) values (?,?)`)
	stmt.Exec(svc, info)
	defer stmt.Close()
}

func main() {
	conf := ParseConfig("config.json")
	db, err := sql.Open("mysql", *conf.Db)
	Error(err)
	defer db.Close()
	InitTable(db)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		svc := r.FormValue("svc")
		info := r.FormValue("info")
		log.Println(svc, info)
		go Insert(db, svc, info)
	})
	http.ListenAndServe("0.0.0.0:9587", nil)
}