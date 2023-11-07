package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	ServiceVersion string
	httpPort       = flag.String("http.port", "8080", "HTTP Port Number")
	MatcherPattern = flag.String("m.pattern", "", "Matcher pattern")
	MatcherType    = flag.String("m.type", "", "Matcher type ()")

	ResponseMatcher *Matcher
)

func main() {
	flag.Parse()

	http.HandleFunc("/", requestHandler)

	if *MatcherPattern != "" && *MatcherType != "" {
		ResponseMatcher = &Matcher{
			Pattern: *MatcherPattern,
			Type:    *MatcherType,
		}
		fmt.Printf("Using response matcher %+v\n", ResponseMatcher)
	}

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

	if ResponseMatcher != nil {
		respFile := ResponseMatcher.Match(r)
		if respFile != "" {
			fmt.Println("Request matcher result:", respFile)

			response, err := getResponse("./responses/" + respFile)
			if err != nil {
				fmt.Printf("Errror getting response from file %s. Error=%v\n", respFile, err)
			} else {
				fmt.Fprint(w, response)
				return
			}
		}
	}

	fmt.Fprint(w, defRes)
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
