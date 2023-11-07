package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// type Config struct {
// 	URLMatcher string
// 	Payload
// }

type Matcher struct {
	// URL: .../key  - no value
	// Query: key1=value1&key2=value2  - possibly URL encoded
	// JSON: { "key": "value" }
	// XML: <main><key>value</key></main>
	Type    string
	Pattern string
}

func (m Matcher) Match(r *http.Request) string {
	data := ""

	switch m.Type {
	case "JSON":
		reqBody, _ := io.ReadAll(r.Body)
		data = string(reqBody)
	case "XML":
		reqBody, _ := io.ReadAll(r.Body)
		data = string(reqBody)
	case "URL":
		data = r.URL.Path
	case "QueryString":
		data = r.URL.RawQuery
	}

	fmt.Println("Matcher using data", data)

	return m.match(data)
}

func (m Matcher) match(dataString string) string {
	var data map[string]string

	switch m.Type {
	case "JSON":
		data = decodeJSON(dataString)
	case "XML":
		data = decodeXML(dataString)
	case "URL":
		data = decodeURL(dataString)
	case "QueryString":
		data = decodeQueryString(dataString)
	}

	if m.Pattern != "" {
		for k, v := range data {
			if match, err := regexp.MatchString(m.Pattern, k); err == nil && match {
				return v
			}
		}
	}

	return ""
}

func decodeJSON(jsonString string) map[string]string {
	var data map[string]string

	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		fmt.Println("JSON error:", err)
	}

	return data
}

func decodeXML(xmlString string) map[string]string {
	key := ""
	data := make(map[string]string)

	p := xml.NewDecoder(strings.NewReader(xmlString))
	for token, err := p.Token(); err == nil; token, err = p.Token() {
		// fmt.Println("XML values", key)
		switch t := token.(type) {
		case xml.StartElement:
			// fmt.Println("Found StartElement elem", t)
			if key == "" {
				// key = opening tag
				key = t.Name.Local
			} else {
				// key = nested opening tags
				key += "." + t.Name.Local
			}
			for _, a := range t.Attr {
				// add attributes to key
				key += fmt.Sprintf("{%s=%s}", a.Name.Local, a.Value)
			}
		case xml.CharData:
			// fmt.Println("Found CharData elem", string([]byte(t)))
			data[key] = string([]byte(t))
		case xml.EndElement:
			// fmt.Println("Found EndElement elem", t)
			if key != "" {
				// remove last opening tag from key
				s := strings.Split(key, ".")
				key = strings.Join(s[:len(s)-1], ".")
			}
		}
	}

	return data
}

func decodeURL(urlString string) map[string]string {
	data := make(map[string]string)
	tokens := strings.Split(urlString, "/")

	for _, t := range tokens {
		data[t] = t
	}

	return data
}

func decodeQueryString(queryString string) map[string]string {
	data := make(map[string]string)

	queryString, e := url.QueryUnescape(queryString)
	if e != nil {
		return data
	}

	if strings.Contains(queryString, "?") {
		split := strings.Split(queryString, "?")
		queryString = split[1]
	}

	kvPairs := strings.Split(queryString, "&")
	for _, kv := range kvPairs {
		if strings.Contains(kv, "=") {
			split := strings.Split(kv, "=")
			data[split[0]] = split[1]
		} else {
			data[kv] = kv
		}
	}

	return data
}
