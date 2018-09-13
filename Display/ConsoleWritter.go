package Display

type ConsoleWritter struct {
	writeFunc func(string) (int, error)
}

func (w *ConsoleWritter) Write(p []byte) (n int, err error) {
	stringData := string(p)
	return w.writeFunc(stringData)
}

func NewConsoleWritter(writeFunc func(string) (int, error)) *ConsoleWritter {
	return &ConsoleWritter{writeFunc: writeFunc}
}
