package msgpack_test

import (
	"math"
	. "msgpackconv/msgpack"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJSONBasicType(t *testing.T) {
	type args struct {
		msgpackconv []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"fixstr",
			args{[]byte{0xa1, 0x61}},
			[]byte(`"a"`),
		},
		{
			"str8",
			args{append([]byte{0xd9, 0x40}, slices.Repeat([]byte{0x61}, int(math.Pow(2, 6)))...)},
			[]byte("\"" + strings.Repeat("a", int(math.Pow(2, 6))) + "\""),
		},
		{
			"str16",
			args{append([]byte{0xda, 0x01, 0x00}, slices.Repeat([]byte{0x61}, int(math.Pow(2, 8)))...)},
			[]byte("\"" + strings.Repeat("a", int(math.Pow(2, 8))) + "\""),
		},
		{
			"str32",
			args{append([]byte{0xdb, 0x00, 0x01, 0x00, 0x00}, slices.Repeat([]byte{0x61}, int(math.Pow(2, 16)))...)},
			[]byte("\"" + strings.Repeat("a", int(math.Pow(2, 16))) + "\""),
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
			"float32",
			args{[]byte(`0.1`)},
			[]byte{0xca, 0x3d, 0xcc, 0xcc, 0xcd},
		},
		{
			"float64",
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
			if got := ToJSON(tt.args.msgpackconv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJSONArray(t *testing.T) {
	getArgs := func(exp int) []byte {
		arr := slices.Repeat([]string{"0"}, int(math.Pow(2, float64(exp))))
		return []byte("[" + strings.Join(arr, ",") + "]")
	}
	type args struct {
		msgpackconv []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"fixarray",
			args{[]byte{0x91, 0x01}},
			[]byte(`[1]`),
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
			if got := ToJSON(tt.args.msgpackconv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJSONMap(t *testing.T) {
	type args struct {
		msgpackconv []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"fixmap",
			args{[]byte{0x86, 0xa3, 0x73, 0x74, 0x72, 0xa1, 0x61, 0xa3, 0x69, 0x6e, 0x74, 0x01, 0xa5, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0xcb, 0x3f, 0xf3, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0xa3, 0x6e, 0x69, 0x6c, 0xc0, 0xa5, 0x61, 0x72, 0x72, 0x61, 0x79, 0x92, 0x00, 0x00, 0xa3, 0x6d, 0x61, 0x70, 0x81, 0xa3, 0x73, 0x74, 0x72, 0xa1, 0x61}},
			[]byte(`{"str": "a","int": 1,"float": 1.2,"nil": null,"array": [0,0],"map": {"str": "a"}}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJSON(tt.args.msgpackconv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
	assert.JSONEq(
		t,
		string(ToJSON([]byte{0xde, 0x00, 0x10, 0xa1, 0x61, 0x01, 0xa1, 0x62, 0x01, 0xa1, 0x63, 0x01, 0xa1, 0x64, 0x01, 0xa1, 0x65, 0x01, 0xa1, 0x66, 0x01, 0xa1, 0x67, 0x01, 0xa1, 0x68, 0x01, 0xa1, 0x69, 0x01, 0xa1, 0x6a, 0x01, 0xa1, 0x6b, 0x01, 0xa1, 0x6c, 0x01, 0xa1, 0x6d, 0x01, 0xa1, 0x6e, 0x01, 0xa1, 0x6f, 0x01, 0xa1, 0x70, 0x01})),
		`{"a":1,"b":1,"c":1,"d":1,"e":1,"f":1,"g":1,"h":1,"i":1,"j":1,"k":1,"l":1,"m":1,"n":1,"o":1,"p":1 }`,
		"test map16 fail",
	)
}

func TestToJSONFail(t *testing.T) {
	type args struct {
		msgpackconv []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"invalid message pack",
			args{[]byte{0x81, 0xa3, 0x69}},
			[]byte{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJSON(tt.args.msgpackconv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
