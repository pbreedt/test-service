package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	ServiceVersion string
	httpPort       = flag.String("http.port", "8080", "HTTP Port Number")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", requestHandler)

	fmt.Printf("Starting service at http://localhost:%s/\n", *httpPort)
	http.ListenAndServe(fmt.Sprintf(":%s", *httpPort), nil)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fmt.Println("Got request URL", r.URL.Path)
	reqBody, _ := io.ReadAll(r.Body)
	fmt.Println("Got request body", string(reqBody))

	defRes, derr := getResponse("./responses/default")
	if derr != nil {
		defRes = `{ "result" : "success" }`
	}
	respFile := ""
	pathArr := strings.Split(r.URL.Path, "/")

	if len(pathArr) >= 2 {
		respFile = pathArr[1]
		fmt.Println("Required response provided in URL", respFile)

		scope := getScope(string(reqBody))
		if scope != "" {
			respFile += "." + scope
		}

		response, err := getResponse("./responses/" + respFile)
		if err != nil {
			fmt.Printf("Errror getting response from file %s. Error=%v\n", respFile, err)
		} else {
			fmt.Fprint(w, response)
			return
		}
	}

	fmt.Fprint(w, defRes)
}

func getScope(body string) string {
	if len(body) <= 0 {
		return ""
	}

	unesc, e := url.QueryUnescape(string(body))
	fmt.Println("Got request body", unesc, "error", e)

	fields := strings.Split(unesc, "&")
	for _, field := range fields {
		kv := strings.Split(field, "=")
		if len(kv) == 2 && kv[0] == "scope" && len(kv[1]) > 0 {
			return kv[1]
		}
	}

	return ""
}

func getResponse(reqPath string) (string, error) {
	resFile, err := os.Open(reqPath)
	if err != nil {
		return "", err
	}

	defer resFile.Close()

	resBytes, err := io.ReadAll(resFile)
	if err != nil {
		return "", err
	}

	fmt.Println("Respond with:", string(resBytes))

	return string(resBytes), nil
}
