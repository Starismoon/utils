package utils

import (
	"encoding/json"
	"fmt"
)

func StructToMapDemo(obj interface{}) map[string]interface{} {
	jsonStr,err:=json.Marshal(obj)
	var mapResult map[string]interface{}
	err = json.Unmarshal(jsonStr, &mapResult)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}
	return mapResult
}
