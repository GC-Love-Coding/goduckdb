package common

// For simplicity, ignore efficiency.
func Checksum(buffer []byte) uint64 {
	return Hash(buffer)
}
