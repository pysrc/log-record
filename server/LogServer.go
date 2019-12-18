package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/http"
)

func main() {
	var host = flag.String("h", "0.0.0.0:9587", "服务地址")
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		svc := r.FormValue("svc")
		info := r.FormValue("info")
		dec, err := base64.StdEncoding.DecodeString(info)
		if err != nil {
			return
		}
		log.Println(svc, string(dec))
	})
	var err = http.ListenAndServe(*host, nil)
	if err != nil {
		log.Panic(err)
	}
}
