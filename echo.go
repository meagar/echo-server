package main

import (
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

func init() {
	echoTemplate = template.Must(template.New("").Parse(echoTemplateStr))
}

var echoTemplate *template.Template
var echoTemplateStr = `
<!doctype html>
<html lang="en">
<head>
	<title>ECHO!</title>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
</head>
<body>
	<div class="container">
		<h1>Echo!</h1>

		<p>This service echos details about each request back to the client.</p>

		<h3>Request Details</h3>

		<table class="table">
		<tr>
			<th>Args</th>
			<td>{{.Args}}</td>
		</tr>
		<tr>
			<th>Method</th>
			<td>{{.Method}}</td>
		</tr>
		<tr>
			<th>URL</th>
			<td>{{.URL}}</td>
		</tr>
		<tr>
			<th>Header</th>
			<td>
				<table>
					{{range $key, $values := .Header}}
						<tr>
							<th>{{$key}}</th>
							<td>
							<ul>
							{{range $value := $values}}
								<li>{{$value}}</li>
							{{end}}
							</ul>
							</td>
						</tr>
					{{end}}
				</table>
			</td>
		</tr>
		<tr>
			<th>ENV</th>
			<td>
				<table>
					{{range $key, $value := .Env}}
						<tr>
							<th>{{$key}}</th>
							<td>{{$value}}</th>
						</tr>
					{{end}}
				</table>
			</td>
		</tr>
		</table>
	</div>
</body>
`
