package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

// Level はログレベル
type Level int

const (
	de Level = iota
	in
	er
)

// Logger は ログレベルとイベントかエラーを受け取ってログに出力する
func Logger(l Level, e string) {
	sl := ""
	switch l {
	case de:
		sl = "[DEBUG]"
	case in:
		sl = "[INFO]"
	case er:
		sl = "[ERROR]"
	default:
		panic("unexpected loglevel")
	}

	entry := map[string]string{
		"severity": sl,
		"message":  e,
		"time":     time.Now().Format(time.RFC3339Nano),
	}
	bytes, _ := json.Marshal(entry)
	fmt.Println(string(bytes))

}
