package main

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/tidwall/gjson"
)

func GetLic(product string, dur int) (isOk bool, result string) {
	req := httplib.Get(host + "/getLic?device=" + getMacMD5() + "&dur=" + fmt.Sprint(dur) + "&product=" + product)
	res, err := req.String()
	if err != nil {
		isOk = false
		result = err.Error()
		return
	}
	code := gjson.Get(res, "code").Int()
	msg := gjson.Get(res, "msg").String()
	result = msg
	if code != 0 {
		isOk = false
		return
	}
	isOk = true
	return
}
