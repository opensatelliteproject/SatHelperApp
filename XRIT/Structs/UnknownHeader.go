package Structs

type UnknownHeader struct {
	Type byte
	Data []byte
}

func MakeUnknownHeader(headerType byte, data []byte) *UnknownHeader {
	v := UnknownHeader{}

	v.Type = headerType

	v.Data = data

	return &v
}

func (uh *UnknownHeader) GetType() int {
	return int(uh.Type)
}
