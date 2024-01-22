package customfunc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// 动态脚本中经常使用的函数封装
//IncreaPageNumber 递增分页参数,实现翻页效果
func IncreaPageNumber(input string) (output string, err error) {
	bodyPath := "body"
	body := gjson.Get(input, fmt.Sprintf("%s.@fromstr", bodyPath)).String()
	pageNumPath := "_param.pageNum"
	index := gjson.Get(body, pageNumPath).Int()
	if index == 0 { // 页码从0开始
		return output, nil
	}
	index++
	body, err = sjson.Set(body, pageNumPath, strconv.Itoa(int(index)))
	if err != nil {
		return "", err
	}
	output, err = sjson.Set(input, bodyPath, body)
	if err != nil {
		return "", err
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
