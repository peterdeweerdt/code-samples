package pjd

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// HTTPClient helps to standardize and encapsulate some common usage patterns
type HTTPClient struct {
	BaseURL     string
	BasicAuth   string
	Logging     bool
	ContentType string
}

// Get will perform an HTTP GET request using the existing base url (if present) and the given path.
// NOTE: the returned http.Response will already have it's body read and you will not be able to re-read.
// Pass in a response body instead if you want access to the body.
func (c HTTPClient) Get(path string, resBody interface{}) (*http.Response, error) {
	return c.request("GET", path, nil, resBody)
}

// Post will perform an HTTP POST request using the existing base url (if present), the given path and request body.
// NOTE: the returned http.Response will already have it's body read and you will not be able to re-read.
// Pass in a response body instead if you want access to the body.
func (c HTTPClient) Post(path string, reqBody, resBody interface{}) (*http.Response, error) {
	return c.request("POST", path, reqBody, resBody)
}

// Put will perform an HTTP PUT request using the existing base url (if present), the given path and request body.
// NOTE: the returned http.Response will already have it's body read and you will not be able to re-read.
// Pass in a response body instead if you want access to the body.
func (c HTTPClient) Put(path string, reqBody, resBody interface{}) (*http.Response, error) {
	return c.request("PUT", path, reqBody, resBody)
}

// Delete will perform an HTTP DELETE request using the existing base url (if present) and the given path
func (c HTTPClient) Delete(path string) (*http.Response, error) {
	return c.request("DELETE", path, nil, nil)
}

func (c HTTPClient) request(verb, path string, reqBody, respBody interface{}) (*http.Response, error) {
	var url string
	if strings.HasPrefix(path, "http") {
		url = path
	} else {
		if path[0] != '/' {
			return nil, errors.New("path must begin with a '/'")
		}
		if path[len(path)-1] == '/' {
			return nil, errors.New("path must not end with a '/'")
		}
		url = fmt.Sprintf("%s%s", c.BaseURL, path)
	}

	var bodyReader io.Reader
	if verb == "POST" || verb == "PUT" {
		b := []byte{}
		var err error
		if c.ContentType == "application/json" {
			b, err = json.Marshal(reqBody)
			if err != nil {
				return nil, err
			}
		} else if c.ContentType == "application/xml" {
			b, err = xml.Marshal(reqBody)
			if err != nil {
				return nil, err
			}
		} else {
			str, ok := reqBody.(string)
			if !ok {
				return nil, errors.New("failed to marshal request body")
			}
			b = []byte(str)
		}

		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(verb, url, bodyReader)
	if err != nil {
		return nil, err
	}
	if len(c.BasicAuth) > 0 {
		req.Header.Add("Authorization", "Basic "+c.BasicAuth)
	}
	req.Header.Add("Content-Type", c.ContentType)

	if c.Logging {
		reqBytes, _ := httputil.DumpRequest(req, true)
		log.Println(string(reqBytes))
	}

	client := http.Client{}

	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		req.Header.Set("Authorization", via[0].Header.Get("Authorization"))
		return nil
	}

	start := time.Now()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Logging {
		log.Printf("Request: %s %s took %s\n", verb, req.URL, time.Since(start))
		resBytes, _ := httputil.DumpResponse(res, true)
		log.Println(string(resBytes))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if len(body) > 0 {
		if c.ContentType == "application/json" {
			err = json.Unmarshal(body, &respBody)
			if err != nil {
				return nil, err
			}
		} else if c.ContentType == "application/xml" {
			err = xml.Unmarshal(body, &respBody)
			if err != nil {
				return nil, err
			}
		}
	}

	if res.StatusCode >= 400 {
		return res, errors.New(res.Status)
	}

	return res, nil
}
