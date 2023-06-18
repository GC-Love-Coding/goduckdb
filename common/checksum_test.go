package common

import "testing"

func TestChecksum(t *testing.T) {
	buf1 := []byte{1, 2, 3}
	buf2 := []byte{1, 2, 3}
	checkSum1 := Checksum(buf1)
	checkSum2 := Checksum(buf2)

	if checkSum1 != checkSum2 {
		t.Errorf("Expect Checksum({1,2,3}) = Checksum({1, 2, 3}), got false")
	}

	buf3 := []byte{1, 2, 2}
	checkSum3 := Checksum(buf3)
	if checkSum1 == checkSum3 {
		t.Errorf("Expect Checksum({1,2,2}) != Checksum({1, 2, 3}), got true")
	}
}
