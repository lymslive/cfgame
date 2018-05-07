package server

import (
	"log"
	"net/http"
	"os/exec"
	"path/filepath"

	// "github.com/lymslive/cfgame/dmlt"
	"github.com/lymslive/cfgame/cmdline"
)

func dealOper(w http.ResponseWriter, r *http.Request) bool {
	script := filepath.Join(cmdline.GetConfig().WebDir, "gopartial.pl")
	pCmd := exec.Command("perl", script)
	if pCmd == nil {
		log.Printf("cannot build execute cmd: perl %s", script)
		return false
	}

	// err := dmlt.FromLocal("")
	out, err := pCmd.Output()
	if err != nil {
		log.Print(err)
		return false
	}

	log.Printf("success execute cmd: perl %s", script)
	// success := []byte("success done!")
	w.Header().Set("Content-Type", "text/plain; charset=GB2312")
	w.Write(out)
	return true
}
