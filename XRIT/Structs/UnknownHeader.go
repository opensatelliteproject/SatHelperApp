package Structs

type UnknownHeader struct {
	Type byte
	Size uint16
	Data []byte
}

func (uh *UnknownHeader) GetType() int {
	return int(uh.Type)
}
