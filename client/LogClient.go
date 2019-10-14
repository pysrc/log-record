package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	svc := flag.String("s", "defalut", "日志系统")
	host := flag.String("h", "http://127.0.0.1:9587", "收集系统推送地址")
	flag.Parse()
	for {
		var log string
		fmt.Scanf("%s\n", &log)
		go http.Get(fmt.Sprintf("%s?svc=%s&info=%s", *host, *svc, log))
	}
}
