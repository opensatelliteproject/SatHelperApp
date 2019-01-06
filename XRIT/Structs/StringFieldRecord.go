package Structs

type StringFieldRecord struct {
	Type     byte
	Size     uint16
	Filename string
}

type AncillaryText StringFieldRecord
type AnnotationRecord StringFieldRecord
type DCSFilenameRecord StringFieldRecord
type HeaderStructuredRecord StringFieldRecord
type ImageDataFunctionRecord StringFieldRecord

func (sfr *StringFieldRecord) GetType() int {
	return int(sfr.Type)
}
func (a *AncillaryText) GetType() int {
	return int(a.Type)
}
func (a *AnnotationRecord) GetType() int {
	return int(a.Type)
}
func (a *DCSFilenameRecord) GetType() int {
	return int(a.Type)
}
func (a *HeaderStructuredRecord) GetType() int {
	return int(a.Type)
}
func (a *ImageDataFunctionRecord) GetType() int {
	return int(a.Type)
}
