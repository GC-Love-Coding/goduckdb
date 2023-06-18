package common

import "fmt"

// efficient hash function that maximizes the avalanche effect and minimizes
// bias
// see: https://nullprogram.com/blog/2018/07/31/
func murmurhash32(x uint32) uint64 {
	x ^= x >> 16
	x *= 0x85ebca6b
	x ^= x >> 13
	x *= 0xc2b2ae35
	x ^= x >> 16

	return uint64(x)
}

func murmurhash64(x uint64) uint64 {
	x ^= x >> 30
	x *= uint64(0xbf58476d1ce4e5b9)
	x ^= x >> 27
	x *= uint64(0x94d049bb133111eb)
	x ^= x >> 31
	return x
}

// TODO: fix me
func hashFloat32(x float32) uint64 {
	panic("hashFloat32.")
}

func hashFloat64(x float64) uint64 {
	panic("hashFloat32.")
}

func hashString(x string) uint64 {
	var hash uint64 = 5381

	for i := range x {
		hash = ((hash << 5) + hash) + uint64(x[i])
	}

	return hash
}

func hashBytes(x []byte) uint64 {
	return hashString(string(x))
}

func Hash(val interface{}) uint64 {
	switch v := val.(type) {
	case uint64:
		return murmurhash64(v)
	case int64:
		return murmurhash64(uint64(v))
	case float32:
		return hashFloat32(v)
	case float64:
		return hashFloat64(v)
	case string:
		return hashString(v)
	case []byte:
		return hashBytes(v)
	default:
		panic(fmt.Sprintf("Default: %T", v))
	}
}
