package common

type Serializer interface {
	WriteData(buffer []byte)
	Write(v interface{})
}
