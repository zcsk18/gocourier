package misc

import (
	"log"
	"strings"
)

func AnalysisHttp(str string) {
	var host string
	var port = "80"

	//log.Printf("first: %s\n", buf)
	if strings.Contains(str, "http") {
		host = str[strings.Index(str, "/")+2:]
	} else {
		host = str[strings.Index(str, " ")+1:]
	}
	if strings.Contains(host, " ") {
		host = host[:strings.Index(host, " ")]
	}
	if strings.Contains(host, "/") {
		host = host[:strings.Index(host, "/")]
	}
	if strings.Contains(host, ":") {
		port = host[strings.Index(host, ":")+1:]
		host = host[:strings.Index(host, ":")]
	}

	log.Printf("host: %s \n", host)
	log.Printf("port: %s \n", port)
}