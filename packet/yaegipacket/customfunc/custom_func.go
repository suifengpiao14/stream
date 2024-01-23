package customfunc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// 动态脚本中经常使用的函数封装

//GetsetJson 指定gjson path确定路径，支持子集多次被序列化情况，修改值,常用来修改翻页参数
func GetsetJson(input string, gjsonPath string, changeFn func(oldValue string) (newValue string, err error)) (output string, err error) {
	if gjsonPath == "" { // 最后一级
		output, err = changeFn(input)
		if err != nil {
			return "", err
		}
		return output, nil
	}
	dotIndex := strings.Index(gjsonPath, ".")
	prePath := gjsonPath
	lastPath := ""
	if dotIndex > -1 {
		prePath, lastPath = gjsonPath[:dotIndex], gjsonPath[dotIndex+1:]
	}
	result := gjson.Get(input, prePath)
	sub, err := GetsetJson(result.String(), lastPath, changeFn)
	if err != nil {
		return "", err
	}
	if result.IsObject() || result.IsArray() {
		output, err = sjson.SetRaw(input, prePath, sub)
	} else {
		output, err = sjson.Set(input, prePath, sub)
	}
	return output, err

}

func WalkArrayMap(body string, gjsonPath string, walkFn func(row map[string]any) (newRow map[string]any)) (output string, err error) {
	data := gjson.Get(body, gjsonPath).String()
	if data == "" {
		return "", nil
	}
	records := make([]map[string]any, 0)
	err = json.Unmarshal([]byte(data), &records)
	if err != nil {
		return "", err
	}

	for i, row := range records {
		records[i] = walkFn(row)
	}
	b, err := json.Marshal(records)
	if err != nil {
		return "", err
	}
	output = string(b)
	return
}

func Fen2yuan(fen any) string {
	var yuan float64
	intFen, ok := fen.(int)
	if ok {
		yuan = float64(intFen) / 100
		return strconv.FormatFloat(yuan, 'f', 2, 64)
	}
	strFen, ok := fen.(string)
	if ok {
		intFen, err := strconv.Atoi(strFen)
		if err == nil {
			yuan = float64(intFen) / 100
			return strconv.FormatFloat(yuan, 'f', 2, 64)
		}
	}
	return strFen
}

func MaskFunc(s string, mark string, start int, length int) (value string) {
	if s == "" {
		return s
	}
	runeS := []rune(s) // 兼容中文
	lastIndex := len(runeS) - 1
	var prefix []rune
	var suffix []rune
	if lastIndex < start {
		prefix = runeS
	}
	if lastIndex >= start {
		prefix = runeS[:start]
	}

	if lastIndex > start+length-1 { //start+length-1  length 包含 start 这个字符
		suffix = runeS[start+length:]
	}

	value = fmt.Sprintf("%s%s%s", string(prefix), mark, string(suffix))
	return value
}
