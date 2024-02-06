package logger

func Example() {
	l, e := NewLogger()
	if e != nil {

	}
	defer l.Shutdown()
	l.Debugln("111")
	l.Warnln("111")

}
