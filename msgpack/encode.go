package msgpack

import (
	"encoding/binary"
	"encoding/json"
	"math"
)

var FirstByte = map[string]byte{
	"positiveFixint": 0x00,
	"uint8":          0xcc,
	"uint16":         0xcd,
	"uint32":         0xce,
	"uint64":         0xcf,
	"negativeFixint": 0xe0,
	"int8":           0xd0,
	"int16":          0xd1,
	"int32":          0xd2,
	"int64":          0xd3,
	"float32":        0xca,
	"float64":        0xcb,
	"fixstr":         0xa0,
	"str8":           0xd9,
	"str16":          0xda,
	"str32":          0xdb,
	"nil":            0xc0,
	"false":          0xc2,
	"true":           0xc3,
	"fixarray":       0x90,
	"array16":        0xdc,
	"array32":        0xdd,
	"fixmap":         0x80,
	"map16":          0xde,
	"map32":          0xdf,
}

func FromJSON(bytes []byte) []byte {
	var obj interface{}
	err := json.Unmarshal(bytes, &obj)
	if err != nil {
		return []byte{}
	}
	ans := []byte{}
	ans = append(ans, encode(obj)...)
	return ans
}

func encode(obj interface{}) []byte {
	ans := []byte{}

	switch v := obj.(type) {
	case nil:
		ans = append(ans, FirstByte["nil"])
	case bool:
		ans = append(ans, getBoolFormat(v))
	case string:
		ans = append(ans, getStrFormat(v)...)
	case float64:
		ans = append(ans, getNumberFormat(v)...)
	case map[string]interface{}:
		ans = append(ans, getMapFormat(v)...)
		for k, vvv := range v {
			ans = append(ans, getStrFormat(k)...)
			ans = append(ans, encode(vvv)...)
		}
	case []interface{}:
		ans = append(ans, getArrayFormat(v)...)
		for i := range v {
			ans = append(ans, encode(v[i])...)
		}
	}
	return ans
}

func getStrFormat(v string) []byte {
	l := len(v)
	var ans []byte
	switch {
	case l < 32:
		// fixstr
		ans = make([]byte, 1+l)
		ans[0] = FirstByte["fixstr"] | byte(len(v))
		copy(ans[1:], []byte(v))
	case float64(l) < math.Pow(2, 8):
		// str8
		ans = make([]byte, 2+l)
		ans[0] = FirstByte["str8"]
		ans[1] = byte(l)
		copy(ans[2:], []byte(v))
	case float64(l) < math.Pow(2, 16):
		// str16
		ans = make([]byte, 3+l)
		ans[0] = FirstByte["str16"]
		binary.BigEndian.PutUint16(ans[1:3], uint16(l))
		copy(ans[3:], []byte(v))
	default:
		// str32
		// because maximum byte size of a String object is (2^32)-1
		ans = make([]byte, 5+l)
		ans[0] = FirstByte["str32"]
		binary.BigEndian.PutUint32(ans[1:5], uint32(l))
		copy(ans[5:], []byte(v))
	}
	return ans
}

func getBoolFormat(v bool) byte {
	if v {
		return FirstByte["true"]
	}
	return FirstByte["false"]
}

func getNumberFormat(v float64) []byte {
	if v == float64(int(v)) {
		if v >= 0 {
			return getPositiveIntFormat(int(v))
		}
		return getNegativeIntFormat(int(v))
	}
	return getFloatFormat(v)
}

func getPositiveIntFormat(v int) []byte {
	switch {
	case v < 128:
		return []byte{FirstByte["positiveFixint"] | byte(v)}
	case float64(v) < math.Pow(2, 8):
		// uint8
		ans := make([]byte, 2)
		ans[0] = FirstByte["uint8"]
		ans[1] = byte(v)
		return ans
	case float64(v) < math.Pow(2, 16):
		// uint16
		ans := make([]byte, 3)
		ans[0] = FirstByte["uint16"]
		binary.BigEndian.PutUint16(ans[1:3], uint16(v))
		return ans
	case float64(v) < math.Pow(2, 32):
		// uint32
		ans := make([]byte, 5)
		ans[0] = FirstByte["uint32"]
		binary.BigEndian.PutUint32(ans[1:5], uint32(v))
		return ans
	default:
		// uint64
		ans := make([]byte, 9)
		ans[0] = FirstByte["uint64"]
		binary.BigEndian.PutUint64(ans[1:9], uint64(v))
		return ans
	}
}

func getNegativeIntFormat(v int) []byte {
	switch {
	case -v <= 32:
		return []byte{FirstByte["negativeFixint"] | byte(32+v)}
	case float64(-v) < math.Pow(2, 7):
		// int8
		ans := make([]byte, 2)
		ans[0] = FirstByte["int8"]
		ans[1] = byte(v)
		return ans
	case float64(-v) < math.Pow(2, 15):
		// int16
		ans := make([]byte, 3)
		ans[0] = FirstByte["int16"]
		binary.BigEndian.PutUint16(ans[1:3], uint16(v))
		return ans
	case float64(-v) < math.Pow(2, 31):
		// int32
		ans := make([]byte, 5)
		ans[0] = FirstByte["int32"]
		binary.BigEndian.PutUint32(ans[1:5], uint32(v))
		return ans
	default:
		// int64
		ans := make([]byte, 9)
		ans[0] = FirstByte["int64"]
		binary.BigEndian.PutUint64(ans[1:9], uint64(v))
		return ans
	}
}

func getFloatFormat(v float64) []byte {
	bits := math.Float64bits(v)
	ans := make([]byte, 9)
	ans[0] = FirstByte["float64"]
	binary.BigEndian.PutUint64(ans[1:9], bits)
	return ans
}

func getMapFormat(v map[string]interface{}) []byte {
	l := len(v)
	switch {
	case l < 16:
		// fixmap
		return []byte{(FirstByte["fixmap"] | byte(len(v)))}
	case l < int(math.Pow(2, 16)):
		// map16
		ans := make([]byte, 3)
		ans[0] = FirstByte["map16"]
		binary.BigEndian.PutUint16(ans[1:3], uint16(l))
		return ans
	default:
		// map32
		ans := make([]byte, 5)
		ans[0] = FirstByte["map32"]
		binary.BigEndian.PutUint32(ans[1:5], uint32(l))
		return ans
	}
}

func getArrayFormat(v []interface{}) []byte {
	l := len(v)
	switch {
	case l < 16:
		// fixarray
		return []byte{(FirstByte["fixarray"] | byte(len(v)))}
	case l < int(math.Pow(2, 16)):
		// array16
		ans := make([]byte, 3)
		ans[0] = FirstByte["array16"]
		binary.BigEndian.PutUint16(ans[1:3], uint16(l))
		return ans
	default:
		// array32
		ans := make([]byte, 5)
		ans[0] = FirstByte["array32"]
		binary.BigEndian.PutUint32(ans[1:5], uint32(l))
		return ans
	}
}
