package main

import (
	"encoding/json"
	"strconv"
)

func code(code int, data interface{}) string {
	codestr := strconv.Itoa(code)
	msg := errcode[codestr]
	msg_str := ""
	if msg.Exists() {
		msg_str = msg.String()
	}
	redata := map[string]interface{}{
		"code": code,
		"msg":  msg_str,
		"data": data,
	}
	bytes, _ := json.Marshal(redata)
	return string(bytes)
}
