package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type Handler struct {
	whitelistMap map[string]bool
}

func initLog(
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	initLog(os.Stdout, os.Stdout, os.Stderr)

	whitelist := os.Getenv("WHITELIST")
	if whitelist == "" {
		fmt.Printf("WHITELIST not set\n")
		os.Exit(1)
	}

	h := &Handler{
		whitelistMap: make(map[string]bool),
	}

	for _, v := range strings.Split(whitelist, ",") {
		Info.Printf("whitelisting: %s\n", v)
		h.whitelistMap[v] = true
	}

	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = true
	proxy.OnRequest().HandleConnectFunc(h.handleConnect)
	proxy.OnRequest().DoFunc(h.handleRequest)

	log.Fatal(http.ListenAndServe(":8080", proxy))
}

func (h *Handler) handleRequest(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	var host string
	if !strings.Contains(r.URL.Host, ":") {
		host = r.URL.Host + ":80"
	} else {
		host = r.URL.Host
	}
	if _, ok := h.whitelistMap[host]; !ok {
		Info.Printf("Host not allowed: %s\n", host)
		return r, goproxy.NewResponse(r,
			goproxy.ContentTypeText, http.StatusForbidden,
			"Host not in whitelist\n")
	}
	return r, nil
}
func (h *Handler) handleConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	if _, ok := h.whitelistMap[host]; !ok {
		Info.Printf("Host not allowed: %s\n", host)
		return goproxy.RejectConnect, host
	}
	return goproxy.OkConnect, host
}
