package pjd

import (
	"log"
	"net/http"
	"net/http/httputil"
)

func MustDumpResponse(res *http.Response) string {
	bytes, err := httputil.DumpResponse(res, true)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error)
	}
	return string(bytes)
}

func MustDumpRequest(res *http.Request) string {
	bytes, err := httputil.DumpRequest(res, true)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error)
	}
	return string(bytes)
}
