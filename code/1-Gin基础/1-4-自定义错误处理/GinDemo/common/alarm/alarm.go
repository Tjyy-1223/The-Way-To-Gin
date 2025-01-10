package alarm

import (
	"GinDemo/common/function"
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type errorString struct {
	s string
}

type errorInfo struct {
	Time     string `json:"time"`
	Alarm    string `json:"alarm"`
	Message  string `json:"message"`
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Funcname string `json:"funcname"`
}

func (e *errorString) Error() string {
	return e.s
}

func New(text string) error {
	alarm("INFO", text, 2)
	return &errorString{text}
}

// Email 发邮件
func Email(text string) error {
	alarm("EMAIL", text, 2)
	return &errorString{text}
}

// Sms 发短信
func Sms(text string) error {
	alarm("SMS", text, 2)
	return &errorString{text}
}

// WeChat 发微信
func WeChat(text string) error {
	alarm("WX", text, 2)
	return &errorString{text}
}

// Panic 异常
func Panic(text string) error {
	alarm("PANIC", text, 5)
	return &errorString{text}
}

func alarm(level string, text string, skip int) {
	currentTime := function.GetTimeStr()

	// 定义 文件名、行号、方法名
	fileName, line, functionName := "", 0, ""

	pc, fileName, line, ok := runtime.Caller(skip)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionName = filepath.Ext(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
	}

	var msg = errorInfo{
		Time:     currentTime,
		Alarm:    level,
		Message:  text,
		Filename: fileName,
		Line:     line,
		Funcname: functionName,
	}

	jsons, errs := json.Marshal(msg)
	if errs != nil {
		fmt.Println("json marshal error:", errs)
	}

	errorJsonInfo := string(jsons)
	fmt.Println(errorJsonInfo)
}
