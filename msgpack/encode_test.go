package msgpack_test

import (
	"fmt"
	"math"
	. "msgpackconv/msgpack"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromJSONBasicType(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"fixstr",
			args{[]byte(`"a"`)},
			[]byte{0xa1, 0x61},
		},
		{
			"str8",
			args{[]byte("\"" + strings.Repeat("a", 32) + "\"")},
			append([]byte{0xd9, 0x20}, slices.Repeat([]byte{0x61}, 32)...),
		},
		{
			"str16",
			args{[]byte("\"" + strings.Repeat("a", int(math.Pow(2, 8))) + "\"")},
			append([]byte{0xda, 0x01, 0x00}, slices.Repeat([]byte{0x61}, int(math.Pow(2, 8)))...),
		},
		{
			"str32",
			args{[]byte("\"" + strings.Repeat("a", int(math.Pow(2, 16))) + "\"")},
			append([]byte{0xdb, 0x00, 0x01, 0x00, 0x00}, slices.Repeat([]byte{0x61}, int(math.Pow(2, 16)))...),
		},
		{
			"positive fixint",
			args{[]byte(`1`)},
			[]byte{0x01},
		},
		{
			"uint8",
			args{[]byte(`128`)},
			[]byte{0xcc, 0x80},
		},
		{
			"uint16",
			args{[]byte(`256`)},
			[]byte{0xcd, 0x01, 0x00},
		},
		{
			"uint32",
			args{[]byte(`65536`)},
			[]byte{0xce, 0x00, 0x01, 0x00, 0x00},
		},
		{
			"uint64",
			args{[]byte(`4294967296`)},
			[]byte{0xcf, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
		},
		{
			"negative fixint",
			args{[]byte(`-1`)},
			[]byte{0xff},
		},
		{
			"int8",
			args{[]byte(`-33`)},
			[]byte{0xd0, 0xdf},
		},
		{
			"int16",
			args{[]byte(`-128`)},
			[]byte{0xd1, 0xff, 0x80},
		},
		{
			"int32",
			args{[]byte(`-32768`)},
			[]byte{0xd2, 0xff, 0xff, 0x80, 0x00},
		},
		{
			"int64",
			args{[]byte(`-2147483648`)},
			[]byte{0xd3, 0xff, 0xff, 0xff, 0xff, 0x80, 0x00, 0x00, 0x00},
		},
		{
			"float",
			args{[]byte(`0.1`)},
			[]byte{0xcb, 0x3f, 0xb9, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a},
		},
		{
			"nil",
			args{[]byte(`null`)},
			[]byte{0xc0},
		},
		{
			"false",
			args{[]byte(`false`)},
			[]byte{0xc2},
		},
		{
			"true",
			args{[]byte(`true`)},
			[]byte{0xc3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromJSON(tt.args.bytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromJSONArray(t *testing.T) {
	getArgs := func(exp int) []byte {
		arr := slices.Repeat([]string{"0"}, int(math.Pow(2, float64(exp))))
		return []byte("[" + strings.Join(arr, ", ") + "]")
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"fixarray",
			args{[]byte(`[]`)},
			[]byte{0x90},
		},
		{
			"array16",
			args{getArgs(4)},
			append([]byte{0xdc, 0x00, 0x10}, slices.Repeat([]byte{0x00}, int(math.Pow(2, 4)))...),
		},
		{
			"array32",
			args{getArgs(16)},
			append([]byte{0xdd, 0x00, 0x01, 0x00, 0x00}, slices.Repeat([]byte{0x00}, int(math.Pow(2, 16)))...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromJSON(tt.args.bytes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromJSONMap(t *testing.T) {
	// 注意 map 的結果不會按照原本的 field 順序
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"str fixmap",
			args{[]byte(`{"str": "a"}`)},
			[]byte{0x81, 0xa3, 0x73, 0x74, 0x72, 0xa1, 0x61},
		},
		{
			"int fixmap",
			args{[]byte(`{"int": 1}`)},
			[]byte{0x81, 0xa3, 0x69, 0x6e, 0x74, 0x01},
		},
		{
			"float fixmap",
			args{[]byte(`{"float": 1.2}`)},
			[]byte{0x81, 0xa5, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0xcb, 0x3f, 0xf3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33},
		},
		{
			"nil fixmap",
			args{[]byte(`{"nil": null}`)},
			[]byte{0x81, 0xa3, 0x6e, 0x69, 0x6c, 0xc0},
		},
		{
			"array fixmap",
			args{[]byte(`{"array": [0,0]}`)},
			[]byte{0x81, 0xa5, 0x61, 0x72, 0x72, 0x61, 0x79, 0x92, 0x00, 0x00},
		},
		{
			"nested fixmap",
			args{[]byte(`{"map": {"str": "a"}}`)},
			[]byte{0x81, 0xa3, 0x6d, 0x61, 0x70, 0x81, 0xa3, 0x73, 0x74, 0x72, 0xa1, 0x61},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromJSON(tt.args.bytes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
	assert.Equal(
		t,
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

func TestFromJSONFail(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"invalid json",
			args{[]byte(`{`)},
			[]byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromJSON(tt.args.bytes)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
