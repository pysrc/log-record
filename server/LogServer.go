package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type LogInfo struct {
	Svc  string
	Info string
}

type LogWriter struct {
	Svc    string
	File   *os.File
	Writer *bufio.Writer
}

var logs chan *LogInfo
var svcs []*LogWriter

func Handle() {
	for {
		ilog := <-logs
		fmt.Println(ilog.Svc, ilog.Info)
		for _, svc := range svcs {
			if ilog.Svc == svc.Svc {
				svc.Writer.WriteString(ilog.Info + "\n")
				svc.Writer.Flush()
			}
		}
	}
}

func AddSvc(svc string) string {
	for _, v := range svcs {
		if v.Svc == svc {
			return ""
		}
	}
	var wt LogWriter
	file, err := os.OpenFile(svc+".log", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err.Error()
	}
	wt.Svc = svc
	wt.File = file
	wt.Writer = bufio.NewWriter(file)
	svcs = append(svcs, &wt)
	return "Success !"
}

func main() {

	var host = flag.String("h", "0.0.0.0:9587", "服务地址")
	var size = flag.Int("s", 100, "日志缓存大小")
	flag.Parse()
	logs = make(chan *LogInfo, *size)
	go Handle()
	defer func() {
		for _, v := range svcs {
			v.Writer.Flush()
			v.File.Close()
		}
	}()
	// 新增一个服务日志输出
	http.HandleFunc("/svc", func(w http.ResponseWriter, r *http.Request) {
		svc := r.FormValue("svc")
		w.Write([]byte(AddSvc(svc)))
	})
	// 接收日志推送
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		svc := r.FormValue("svc")
		info := r.FormValue("info")
		dec, err := base64.StdEncoding.DecodeString(info)
		if err != nil {
			return
		}
		logs <- &LogInfo{svc, string(dec)}
	})
	var err = http.ListenAndServe(*host, nil)
	if err != nil {
		panic(err)
	}
}
