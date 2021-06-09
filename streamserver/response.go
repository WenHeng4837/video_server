package main

import (
	"io"
	"net/http"
)

//streamserver的错误返回
func sendErrorResponse(w http.ResponseWriter, sc int, errMsg string) {
	w.WriteHeader(sc)
	io.WriteString(w, errMsg)
}
