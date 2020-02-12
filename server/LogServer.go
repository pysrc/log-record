package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type LogInfo struct {
	Svc  string
	Sid  string
	Info string
}

type LogWriter struct {
	File   *os.File
	Writer *bufio.Writer
}

var logs chan *LogInfo

var svcm map[string]*LogWriter

func Handle() {
	for {
		ilog := <-logs
		fmt.Printf("%s %s %s", ilog.Svc, ilog.Sid, ilog.Info)
		svc := svcm[ilog.Svc]
		if svc != nil {
			_, err := svc.Writer.WriteString(fmt.Sprintf("%s %s", ilog.Sid, ilog.Info))
			if err != nil {
				fmt.Println("system ", err.Error())
				delete(svcm, ilog.Svc)
				AddSvc(ilog.Svc)
				return
			}
			svc.Writer.Flush()
		}
	}
}

func AddSvc(svc string) string {
	if svcm[svc] != nil {
		// 已经存在
		return "The service is logging"
	}
	var wt LogWriter
	file, err := os.OpenFile(svc+".log", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err.Error()
	}
	wt.File = file
	wt.Writer = bufio.NewWriter(file)
	svcm[svc] = &wt
	return "Success !"
}

func main() {

	var host = flag.String("h", "0.0.0.0:9587", "服务地址")
	var size = flag.Int("s", 100, "日志缓存大小")
	flag.Parse()
	logs = make(chan *LogInfo, *size)
	svcm = make(map[string]*LogWriter)
	go Handle()
	defer func() {
		for _, v := range svcm {
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
		sid := r.FormValue("sid")
		logs <- &LogInfo{svc, sid, info}
	})
	var err = http.ListenAndServe(*host, nil)
	if err != nil {
		panic(err)
	}
}
