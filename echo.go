package main

import (
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/urfave/negroni"
	"k8s.io/klog/v2"
)

var statusPath string

func main() {
	flag.StringVar(&statusPath, "status", "", "Where to mount the status endpoint, which simply outputs an HTTP 200 with body \"ok\"")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage: echo <port>")
		os.Exit(1)
	}

	klog.Infof("Listening on port %s", flag.Arg(0))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if statusPath != "" && req.URL.Path == statusPath {
			handleStatus(res)
			return
		}

		fmt.Println(req.Header)
		env := map[string]string{}
		for _, value := range os.Environ() {
			pair := strings.Split(value, "=")
			if len(pair) != 2 {
				panic("no!")
			}
			env[pair[0]] = pair[1]
		}

		data := struct {
			Args   []string
			Method string
			URL    *url.URL
			Header http.Header
			Env    map[string]string
		}{
			Args:   os.Args,
			Method: req.Method,
			URL:    req.URL,
			Header: req.Header,
			Env:    env,
		}
		echoTemplate.Execute(res, data)
	})

	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	klog.Error(http.ListenAndServe(fmt.Sprintf(":%s", flag.Arg(0)), n))
}

func handleStatus(rw http.ResponseWriter) {
	rw.Write([]byte("ok"))
}

//go:embed index.gohtml
var echoTemplateStr string

var echoTemplate *template.Template = template.Must(template.New("").Parse(echoTemplateStr))
