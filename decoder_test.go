package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/rwcarlsen/goexif/tiff"
	"testing"

	"github.com/stretchr/testify/assert"
)

type input struct {
	tgId   string
	tpe    string
	nVals  string
	offset string
	val    string
}

const (
	InputDTByte      = "0001"
	InputDTAscii     = "0002"
	InputDTInt32     = "0003"
	InputDTInt64     = "0004"
	InputDTRational  = "0005"
	InputDTSInt32    = "0008"
	InputDTSByte     = "0006"
	InputDTUndefined = "0007"
	InputDTSShort    = "0008"
	InputDTSLong     = "0009"
	InputDTSRational = "000A"
	InputDTFloat     = "000B"
	InputDTDouble    = "000C"
)
const (
	InputValuesCount2 = "00000002"
)

func Test_single_32bit_integer(t *testing.T) {
	w := decoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(
		input{"0001", InputDTInt32, "00000002", "000A" + "0000", ""},
		binary.BigEndian,
	)
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("one", tag)
	assert.NoError(t, err)
	assert.Equal(t, 10, w.data["one"])
}
func Test_single_r(t *testing.T) {
	w := decoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(
		input{"0001", InputDTRational, "00000001", "00000010", "00000019" + "00000064"},
		binary.BigEndian,
	)
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("one", tag)
	assert.NoError(t, err)
	assert.Equal(t, 0.25, w.data["one"])
}
func buildInput(in input, order binary.ByteOrder) []byte {
	data := make([]byte, 0)
	d, _ := hex.DecodeString(in.tgId)
	data = append(data, d...)
	d, _ = hex.DecodeString(in.tpe)
	data = append(data, d...)
	d, _ = hex.DecodeString(in.nVals)
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
