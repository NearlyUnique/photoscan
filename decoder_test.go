package main

import (
	"bytes"
	"encoding/binary"
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
	w := DecoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(InputDTInt32, 10, "")
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("one", tag)
	assert.NoError(t, err)
	assert.Equal(t, 10, w.data["one"])
}
func Test_single_rational(t *testing.T) {
	w := DecoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(InputDTRational, 1, I4(25, 100))
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("one", tag)
	assert.NoError(t, err)
	assert.Equal(t, 0.25, w.data["one"])
}
func Test_multiple_rational(t *testing.T) {
	w := DecoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(InputDTRational, 3, I4(25, 100)+I4(50, 1000)+I4(75, 10))
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("one", tag)
	assert.NoError(t, err)
	assert.Equal(t, []float64{0.25, 0.05, 7.5}, w.data["one"])
}
func Test_single_string(t *testing.T) {
	w := DecoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(InputDTAscii, 6, "abcdef")
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("one", tag)
	assert.NoError(t, err)
	assert.Equal(t, "abcdef", w.data["one"])
}

func Test_single_DateTime(t *testing.T) {
	w := DecoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(InputDTAscii, 0, "2019:10:11 19:05:41")
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("DateTime", tag)
	assert.NoError(t, err)
	assert.Equal(t, "2019-10-11 19:05:41", w.data["DateTime"])
}

func Test_single_GPS(t *testing.T) {
	w := DecoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(InputDTRational, 3, I4(
		52, 1,
		12, 1,
		30000, 987))
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("GPSLatitude", tag)
	assert.NoError(t, err)
	assert.Equal(t, "52°12'30.3951\"", w.data["GPSLatitude"])
}
func Test_single_GPS_Time(t *testing.T) {
	w := DecoderWalker{}
	w.data = make(map[string]any)

	data := buildInput(InputDTRational, 3, I4(
		18, 1,
		5, 1,
		33, 1))
	raw := bytes.NewReader(data)
	tag, err := tiff.DecodeTag(raw, binary.BigEndian)
	assert.NoError(t, err)
	err = w.Walk("GPSTimeStamp", tag)
	assert.NoError(t, err)
	assert.Equal(t, "18:05:33", w.data["GPSTimeStamp"])
}
