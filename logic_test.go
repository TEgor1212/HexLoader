package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreeting(t *testing.T) {
	assert.Equal(t, "Hello, Alice!", Greeting("Alice"))
	assert.Equal(t, "Hello, Bob!", Greeting("Bob"))
	assert.Equal(t, "Hello, !", Greeting(""))
}

func TestSumRange(t *testing.T) {
	tests := []struct {
		start    int
		end      int
		expected int
	}{
		{1, 5, 15},
		{1, 1, 1},
		{1, 0, 0},
		{0, 3, 6},
		{-2, 2, 0},
		{10, 10, 10},
	}
	for _, tc := range tests {
		got := SumRange(tc.start, tc.end)
		assert.Equal(t, tc.expected, got, "SumRange(%d, %d)", tc.start, tc.end)
	}
}

func TestCalculateChecksum(t *testing.T) {
	tests := []struct {
		data     []byte
		expected byte
	}{
		{[]byte{}, 0},
		{[]byte{0x00}, 0x00},
		{[]byte{0xFF}, 0xFF},
		{[]byte{0x01, 0x02}, 0x03},
		{[]byte{0x10, 0x20, 0x30, 0x40}, 0x40}, // 10 ^ 20 ^ 30 ^ 40 = 40
		{[]byte{0xAB, 0xCD, 0xEF}, 0xAB ^ 0xCD ^ 0xEF},
	}
	for _, tc := range tests {
		got := CalculateChecksum(tc.data)
		assert.Equal(t, tc.expected, got, "CalculateChecksum(%v)", tc.data)
	}
}

func TestHexToBytes(t *testing.T) {
	tests := []struct {
		input    string
		expected []byte
		err      bool
	}{
		{"", []byte{}, false},
		{"00", []byte{0x00}, false},
		{"FF", []byte{0xFF}, false},
		{"0102", []byte{0x01, 0x02}, false},
		{"deadbeef", []byte{0xde, 0xad, 0xbe, 0xef}, false},
		{"DEADBEEF", []byte{0xde, 0xad, 0xbe, 0xef}, false},
		{"a", nil, true},    // odd length
		{"xyz", nil, true},  // invalid characters
		{"0x01", nil, true}, // 0x prefix not supported in this raw parser
		{"1G", nil, true},   // invalid character 'G'
	}
	for _, tc := range tests {
		result, err := HexToBytes(tc.input)
		if tc.err {
			assert.Error(t, err, "HexToBytes(%q)", tc.input)
		} else {
			assert.NoError(t, err, "HexToBytes(%q)", tc.input)
			assert.Equal(t, tc.expected, result, "HexToBytes(%q)", tc.input)
		}
	}
}

func TestHexToBytesRoundTrip(t *testing.T) {
	original := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F} // "Hello"
	hex := "48656C6C6F"
	result, err := HexToBytes(hex)
	assert.NoError(t, err)
	assert.Equal(t, original, result)
}
