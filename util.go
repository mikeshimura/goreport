package goreport

import (
	"os"
	"strings"
)

func AddComma(s string) string {
	if strings.Index(s, ".") == -1 {
		return addCommaSub(s)
	}
	ss := strings.Split(s, ".")
	ss[0] = addCommaSub(ss[0])
	return ss[0] + "." + ss[1]
}

func addCommaSub(s string) string {
	res := ""
	if len(s) < 4 {
		return s
	}
	pos := len(s) % 3
	if pos > 0 {
		res += s[0:pos] + ","
	}
	for i := pos; i < len(s); i += 3 {
		res += s[i : i+3]
		//fmt.Printf("pos %v \n", i)
		if i < len(s)-3 {
			res += ","
		}
	}
	return res
}
func ReadTextFile(filename string, colno int) []interface{} {
	res, _ := os.ReadFile(filename)
	return ReadBytes(res, colno)
}

func ReadBytes(src []byte, colno int) []interface{} {
	lines := strings.Split(string(src), "\n")
	list := make([]interface{}, 0, 100)
	for _, line := range lines {
		line = strings.Replace(line, "\r", "", -1)
		cols := strings.Split(line, "\t")
		if len(cols) < colno {
			continue
		}
		list = append(list, cols)
	}
	return list
}
