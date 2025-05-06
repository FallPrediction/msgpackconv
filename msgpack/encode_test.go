package msgpack_test

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	. "msgpackconv/msgpack"
)

func TestStringFromJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`"a"`))), "a161", "test fixstr fail")
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte("\""+strings.Repeat("a", 32)+"\""))),
		"d920"+strings.Repeat("61", 32),
		"test str8 fail",
	)
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte("\""+strings.Repeat("a", int(math.Pow(2, 8)))+"\""))),
		"da0100"+strings.Repeat("61", int(math.Pow(2, 8))),
		"test str16 fail",
	)
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte("\""+strings.Repeat("a", int(math.Pow(2, 16)))+"\""))),
		"db00010000"+strings.Repeat("61", int(math.Pow(2, 16))),
		"test str32 fail",
	)
}

func TestIntegerFromJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`1`))), "01", "test positive fixint fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`128`))), "cc80", "test uint8 fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`256`))), "cd0100", "test uint16 fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`65536`))), "ce00010000", "test uint32 fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`4294967296`))), "cf0000000100000000", "test uint64 fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`-1`))), "ff", "test negative fixint fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`-33`))), "d0df", "test int8 fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`-128`))), "d1ff80", "test int16 fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`-32768`))), "d2ffff8000", "test int32 fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`-2147483648`))), "d3ffffffff80000000", "test int64 fail")
}

func TestFloatFromJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`0.1`))), "cb3fb999999999999a", "test float fail")
}

func TestNilFromJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`null`))), "c0", "test nil fail")
}

func TestBoolFromJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`false`))), "c2", "test false fail")
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`true`))), "c3", "test true fail")
}

func TestArrayFromJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`[]`))), "90", "test fixarray fail")
	arr := slices.Repeat([]string{"0"}, int(math.Pow(2, 4)))
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte("["+strings.Join(arr, ", ")+"]"))),
		"dc0010"+strings.Repeat("00", int(math.Pow(2, 4))),
		"test array16 fail",
	)
	arr = slices.Repeat([]string{"0"}, int(math.Pow(2, 16)))
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte("["+strings.Join(arr, ", ")+"]"))),
		"dd00010000"+strings.Repeat("00", int(math.Pow(2, 16))),
		"test array32 fail",
	)
}

func TestMapFromJSON(t *testing.T) {
	assert := assert.New(t)
	// 注意 map 的結果不會按照順序
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte(`{"str": "a"}`))),
		"81a3737472a161",
		"test fixmap fail",
	)
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte(`{"int": 1}`))),
		"81a3696e7401",
		"test fixmap fail",
	)
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte(`{"float": 1.2}`))),
		"81a5666c6f6174cb3ff3333333333333",
		"test fixmap fail",
	)
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte(`{"nil": null}`))),
		"81a36e696cc0",
		"test fixmap fail",
	)
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte(`{"array": [0,0]}`))),
		"81a56172726179920000",
		"test fixmap fail",
	)
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte(`{"map": {"str": "a"}}`))),
		"81a36d617081a3737472a161",
		"test fixmap fail",
	)
	assert.Equal(
		fmt.Sprintf("%x", FromJSON([]byte(`{
			"a": 1,
			"b": 1,
			"c": 1,
			"d": 1,
			"e": 1,
			"f": 1,
			"g": 1,
			"h": 1,
			"i": 1,
			"j": 1,
			"k": 1,
			"l": 1,
			"m": 1,
			"n": 1,
			"o": 1,
			"p": 1
		}`)))[:6],
		"de0010",
		"test map16 fail",
	)
}

func TestEncodeFailFromJSON(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(fmt.Sprintf("%x", FromJSON([]byte(`{`))), "", "test encode fail")
}
