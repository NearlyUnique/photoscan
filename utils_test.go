package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

// I2 2 byte int
func I2(v int) string {
	return fmt.Sprintf("%04X", v)
}

// I2 2 byte int in 4 bytes, big endian
func I2_(v int) string {
	return fmt.Sprintf("%04X0000", v)
}

// I1 1 byte int in 4 bytes, big endian
func I1_(v int) string {
	return fmt.Sprintf("%02X000000", v)
}

// I4 4 byte int
func I4(nums ...int) string {
	r := ""
	for _, v := range nums {
		r += fmt.Sprintf("%08X", v)
	}
	return r
}
func Ascii(str string) string {
	r := ""
	for _, v := range str {
		r += fmt.Sprintf("%02X", v)
	}
	return r
}

func buildInput(typeStr string, valueCount int, dataStr string) []byte {
	lengthStr := I4(4 + len(dataStr)) // bytes used to store length
	valueOrCount := I4(valueCount)
	switch typeStr {
	case InputDTAscii:
		valueOrCount = I4(len(dataStr))
		lengthStr = I4(4 + len(dataStr) + 2) // bytes used to store length
		dataStr = Ascii(dataStr)
	case InputDTLong:
		lengthStr = I4(len(dataStr))
	case InputDTByte, InputDTSByte,
		InputDTShort:
		// if the dataStr fits in 32 bits then it goes in the nVal
		switch typeStr {
		case InputDTByte, InputDTSByte:
			lengthStr = I1_(valueCount)
		case InputDTShort:
			lengthStr = I2_(valueCount)
		}
		valueOrCount = I4(2)
		dataStr = ""
	}
	//input   {"0001", "0004", "00000002", "00000010", "1112131415161718"},
	in := input{I2(1), typeStr, valueOrCount, lengthStr, dataStr}
	order := binary.BigEndian

	data := make([]byte, 0)
	d, _ := hex.DecodeString(in.tagID)
	data = append(data, d...)
	d, _ = hex.DecodeString(in.dataType)
	data = append(data, d...)
	d, _ = hex.DecodeString(in.valueCount)
	data = append(data, d...)
	d, _ = hex.DecodeString(in.offset)
	data = append(data, d...)

	if in.val != "" {
		off := order.Uint32(d)
		for i := 0; i < int(off)-12; i++ {
			data = append(data, 0xFF)
		}

		d, _ = hex.DecodeString(in.val)
		data = append(data, d...)
	}

	return data
}
