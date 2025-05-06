package msgpack_test

import (
	"math"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	. "msgpackconv/msgpack"
)

func TestStringToJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(ToJSON([]byte{0xa1, 0x61}), []byte(`"a"`), "test fixstr fail")
	assert.Equal(
		ToJSON(append([]byte{0xd9, 0x40}, slices.Repeat([]byte{0x61}, int(math.Pow(2, 6)))...)),
		[]byte("\""+strings.Repeat("a", int(math.Pow(2, 6)))+"\""),
		"test str8 fail",
	)
	assert.Equal(
		ToJSON(append([]byte{0xda, 0x01, 0x00}, slices.Repeat([]byte{0x61}, int(math.Pow(2, 8)))...)),
		[]byte("\""+strings.Repeat("a", int(math.Pow(2, 8)))+"\""),
		"test str16 fail",
	)
	assert.Equal(
		ToJSON(append([]byte{0xdb, 0x00, 0x01, 0x00, 0x00}, slices.Repeat([]byte{0x61}, int(math.Pow(2, 16)))...)),
		[]byte("\""+strings.Repeat("a", int(math.Pow(2, 16)))+"\""),
		"test str32 fail",
	)
}

func TestIntegerToJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(ToJSON([]byte{0x01}), []byte(`1`), "test positive fixint fail")
	assert.Equal(ToJSON([]byte{0xcc, 0x80}), []byte(`128`), "test uint8 fail")
	assert.Equal(ToJSON([]byte{0xcd, 0x01, 0x00}), []byte(`256`), "test uint16 fail")
	assert.Equal(ToJSON([]byte{0xce, 0x00, 0x01, 0x00, 0x00}), []byte(`65536`), "test uint32 fail")
	assert.Equal(ToJSON([]byte{0xcf, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00}), []byte(`4294967296`), "test uint64 fail")
	assert.Equal(ToJSON([]byte{0xff}), []byte(`-1`), "test negative fixint fail")
	assert.Equal(ToJSON([]byte{0xd0, 0xdf}), []byte(`-33`), "test int8 fail")
	assert.Equal(ToJSON([]byte{0xd1, 0xff, 0x80}), []byte(`-128`), "test int16 fail")
	assert.Equal(ToJSON([]byte{0xd2, 0xff, 0xff, 0x80, 0x00}), []byte(`-32768`), "test int32 fail")
	assert.Equal(ToJSON([]byte{0xd3, 0xff, 0xff, 0xff, 0xff, 0x80, 0x00, 0x00, 0x00}), []byte(`-2147483648`), "test int64 fail")
}

func TestFloatToJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(ToJSON([]byte{0xca, 0x3d, 0xcc, 0xcc, 0xcd}), []byte(`0.1`), "test float32 fail")
	assert.Equal(ToJSON([]byte{0xcb, 0x3f, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}), []byte(`0.1`), "test float64 fail")
}

func TestNilToJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(ToJSON([]byte{0xc0}), []byte(`null`), "test nil fail")
}

func TestBoolToJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(ToJSON([]byte{0xc2}), []byte(`false`), "test false fail")
	assert.Equal(ToJSON([]byte{0xc3}), []byte(`true`), "test true fail")
}

func TestArrayToJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(ToJSON([]byte{0x91, 0x01}), []byte(`[1]`), "test fixarray fail")
	assert.Equal(ToJSON([]byte{0x92, 0x01, 0x01}), []byte(`[1,1]`), "test fixarray fail")
	arr := slices.Repeat([]string{"0"}, int(math.Pow(2, 4)))
	assert.Equal(
		string(ToJSON(append([]byte{0xdc, 0x00, 0x10}, slices.Repeat([]byte{0x00}, int(math.Pow(2, 4)))...))),
		"["+strings.Join(arr, ",")+"]",
		"test array16 fail",
	)
	arr = slices.Repeat([]string{"0"}, int(math.Pow(2, 16)))
	assert.Equal(
		string(ToJSON(append([]byte{0xdd, 0x00, 0x01, 0x00, 0x00}, slices.Repeat([]byte{0x00}, int(math.Pow(2, 16)))...))),
		"["+strings.Join(arr, ",")+"]",
		"test array32 fail",
	)
}

func TestMapToJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(ToJSON([]byte{0x81, 0xa3, 0x69, 0x6e, 0x74, 0x01}), []byte(`{"int":1}`), "test fixmap fail")
	assert.Equal(ToJSON([]byte{0x82, 0xa3, 0x69, 0x6e, 0x74, 0x01, 0xa3, 0x73, 0x74, 0x72, 0xa1, 0x61}), []byte(`{"int":1,"str":"a"}`), "test fixmap fail")
	assert.Equal(
		ToJSON([]byte{0x81, 0xa3, 0x6d, 0x61, 0x70, 0x81, 0xa3, 0x73, 0x74, 0x72, 0xa1, 0x61}),
		[]byte(`{"map":{"str":"a"}}`),
		"test fixmap fail",
	)
	assert.JSONEq(
		string(ToJSON([]byte{0xde, 0x00, 0x10, 0xa1, 0x61, 0x01, 0xa1, 0x62, 0x01, 0xa1, 0x63, 0x01, 0xa1, 0x64, 0x01, 0xa1, 0x65, 0x01, 0xa1, 0x66, 0x01, 0xa1, 0x67, 0x01, 0xa1, 0x68, 0x01, 0xa1, 0x69, 0x01, 0xa1, 0x6a, 0x01, 0xa1, 0x6b, 0x01, 0xa1, 0x6c, 0x01, 0xa1, 0x6d, 0x01, 0xa1, 0x6e, 0x01, 0xa1, 0x6f, 0x01, 0xa1, 0x70, 0x01})),
		`{"a":1,"b":1,"c":1,"d":1,"e":1,"f":1,"g":1,"h":1,"i":1,"j":1,"k":1,"l":1,"m":1,"n":1,"o":1,"p":1 }`,
		"test map16 fail",
	)
}

func TestToJSONFail(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(ToJSON([]byte{0x81, 0xa3, 0x69}), []byte{}, "test to json fail")
}
