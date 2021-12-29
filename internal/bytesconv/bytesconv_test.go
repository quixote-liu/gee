package bytesconv

import (
	"bytes"
	"crypto/rand"
	"testing"
)

var testString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
var testBytes = []byte(testString)

func rawStringToBytes(s string) []byte {
	return []byte(s)
}

func rawBytesToString(b []byte) string {
	return string(b)
}

func TestBytesToString(t *testing.T) {
	data := make([]byte, 1024)
	for i := 0; i < 100; i++ {
		rand.Read(data)
		if rawBytesToString(data) != BytesToString(data) {
			t.Fatalf("don't match")
		}
	}
}

func TestStringToBytes(t *testing.T) {
	if !bytes.Equal(rawStringToBytes(testString), StringToBytes(testString)) {
		t.Fatalf("don't match")
	}
}

// go test -v -run=bytesconv.go -bench="."

func BenchmarkRawBytesToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawBytesToString(testBytes)
	}
}

func BenchmarkRawStringToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rawStringToBytes(testString)
	}
}

func BenchmarkBytesToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BytesToString(testBytes)
	}
}

func BenchmarkStringToBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		StringToBytes(testString)
	}
}
