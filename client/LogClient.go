package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	svc := flag.String("s", "defalut", "日志系统")
	host := flag.String("h", "http://127.0.0.1:9587", "收集系统推送地址")
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		bts := scanner.Bytes()
		fmt.Println(string(bts))
		http.Get(fmt.Sprintf(
			"%s?svc=%s&info=%s",
			*host,
			*svc,
			base64.StdEncoding.EncodeToString(bts),
		))
	}
}
