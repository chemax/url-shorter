package logger

import "fmt"

func Example() {
	l, err := NewLogger()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	l.Debugln("debug")
	l.Warnln("warn")
	l.Infoln("Info")
}
