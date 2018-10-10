package pjd

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func ConvertBodyToMap(body io.Reader) map[string]interface{} {
	bytes, _ := ioutil.ReadAll(body)
	var m map[string]interface{}
	json.Unmarshal(bytes, &m)
	return m
}

func ConvertBodyToSlice(body io.Reader) []map[string]interface{} {
	bytes, _ := ioutil.ReadAll(body)
	var m []map[string]interface{}
	json.Unmarshal(bytes, &m)
	return m
}

func ConvertByteSliceToMap(bytes []byte) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(bytes, &m)
	return m
}
