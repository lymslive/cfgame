package server

import (
	"log"
	"net/http"

	"github.com/lymslive/cfgame/dmlt"
)

func dealOper(w http.ResponseWriter, r *http.Request) bool {
	err := dmlt.FromLocal("")
	if err != nil {
		log.Print(err)
		return false
	}

	success := []byte("success done!")
	w.Write(success)
	return true
}
