package common

type Deserializer interface {
	ReadData(buffer []byte)
	Read(v interface{})
}
