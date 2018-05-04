package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lymslive/cfgame/cmdline"
	"github.com/lymslive/cfgame/dmlt"
	"github.com/lymslive/cfgame/fsop"
)

// 启动服务器
func Start() {
	cfg := cmdline.GetConfig()
	datadir := fsop.Parent(cfg.WebDir)
	dmlt.SetDataDir(datadir)

	setHandler()
	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(http.ListenAndServe(address, nil))
}

func setHandler() {
	http.HandleFunc("/test/", handleTest)
	http.HandleFunc("/oper", handleOper)
	http.HandleFunc("/", handleRoot)
}

// 根回调函数入口
func handleRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("req from: %q, res by: %s", r.URL.Path, "handleRoot")

	/*
		switch {
		case responNote(w, r):
			log.Printf("success to responNote()")
		default:
			log.Printf("fails to respond to: %q", r.URL.Path)
			http.NotFound(w, r)
		}
	*/
}

// 各注册分配回调
func handleTest(w http.ResponseWriter, r *http.Request) {
	log.Printf("req from: %q, res by: %s", r.URL.Path, "handleTest")

	fmt.Fprintf(w, "RUL.Path = %q\n", r.URL.Path)
}

func handleOper(w http.ResponseWriter, r *http.Request) {
	log.Printf("req from: %q, res by: %s", r.URL.Path, "handleOper")

	switch {
	case dealOper(w, r):
		log.Printf("success to dealOper()")
	default:
		log.Printf("fails to respond to: %q", r.URL.Path)
		http.NotFound(w, r)
	}
}
