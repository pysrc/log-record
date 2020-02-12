package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	svc := flag.String("s", "svc", "服务名称")
	sid := flag.String("i", "svc-id", "服务ID,相同服务之间靠ID区分")
	host := flag.String("h", "http://127.0.0.1:9587", "收集系统推送地址")
	flag.Parse()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		bts := scanner.Bytes()
		var lg = string(bts) + "\n"
		fmt.Print(lg)
		http.Post(*host,
			"application/x-www-form-urlencoded",
			strings.NewReader(fmt.Sprintf(`svc=%s&sid=%s&info=%s`, *svc, *sid, lg)),
		)
	}
}
