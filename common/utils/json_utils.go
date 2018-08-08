package utils


import (
	"encoding/json"
)



func Object2Json(v interface{})  string{
	data, err := json.Marshal(v)
	if err==nil {
		return string(data)
	}
	return ""
}

func Json2Object(data string, v interface{}) interface{} {
	err := json.Unmarshal([]byte(data), v)
	if err == nil {
		return v
	}
	return nil
}

