package msgpack

import (
	"encoding/binary"
	"encoding/json"
	"math"
)

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
		ans = append(ans, 0xc0)
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
		ans[0] = 0xa0 | byte(len(v))
		copy(ans[1:], []byte(v))
	case float64(l) < math.Pow(2, 8):
		// str8
		ans = make([]byte, 2+l)
		ans[0] = 0xd9
		ans[1] = byte(l)
		copy(ans[2:], []byte(v))
	case float64(l) < math.Pow(2, 16):
		// str16
		ans = make([]byte, 3+l)
		ans[0] = 0xda
		binary.BigEndian.PutUint16(ans[1:3], uint16(l))
		copy(ans[3:], []byte(v))
	default:
		// str32
		// because maximum byte size of a String object is (2^32)-1
		ans = make([]byte, 5+l)
		ans[0] = 0xdb
		binary.BigEndian.PutUint32(ans[1:5], uint32(l))
		copy(ans[5:], []byte(v))
	}
	return ans
}

func getBoolFormat(v bool) byte {
	if v {
		return 0xc3
	}
	return 0xc2
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
		return []byte{0x00 | byte(v)}
	case float64(v) < math.Pow(2, 8):
		// uint8
		ans := make([]byte, 2)
		ans[0] = 0xcc
		ans[1] = byte(v)
		return ans
	case float64(v) < math.Pow(2, 16):
		// uint16
		ans := make([]byte, 3)
		ans[0] = 0xcd
		binary.BigEndian.PutUint16(ans[1:3], uint16(v))
		return ans
	case float64(v) < math.Pow(2, 32):
		// uint32
		ans := make([]byte, 5)
		ans[0] = 0xce
		binary.BigEndian.PutUint32(ans[1:5], uint32(v))
		return ans
	default:
		// uint64
		ans := make([]byte, 9)
		ans[0] = 0xcf
		binary.BigEndian.PutUint64(ans[1:9], uint64(v))
		return ans
	}
}

func getNegativeIntFormat(v int) []byte {
	switch {
	case -v <= 32:
		return []byte{0xe0 | byte(32+v)}
	case float64(-v) < math.Pow(2, 7):
		// int8
		ans := make([]byte, 2)
		ans[0] = 0xd0
		ans[1] = byte(v)
		return ans
	case float64(-v) < math.Pow(2, 15):
		// int16
		ans := make([]byte, 3)
		ans[0] = 0xd1
		binary.BigEndian.PutUint16(ans[1:3], uint16(v))
		return ans
	case float64(-v) < math.Pow(2, 31):
		// int32
		ans := make([]byte, 5)
		ans[0] = 0xd2
		binary.BigEndian.PutUint32(ans[1:5], uint32(v))
		return ans
	default:
		// int64
		ans := make([]byte, 9)
		ans[0] = 0xd3
		binary.BigEndian.PutUint64(ans[1:9], uint64(v))
		return ans
	}
}

func getFloatFormat(v float64) []byte {
	bits := math.Float64bits(v)
	ans := make([]byte, 9)
	ans[0] = 0xcb
	binary.BigEndian.PutUint64(ans[1:9], bits)
	return ans
}

func getMapFormat(v map[string]interface{}) []byte {
	l := len(v)
	switch {
	case l < 16:
		// fixmap
		return []byte{(0x80 | byte(len(v)))}
	case l < int(math.Pow(2, 16)):
		// map16
		ans := make([]byte, 3)
		ans[0] = 0xde
		binary.BigEndian.PutUint16(ans[1:3], uint16(l))
		return ans
	default:
		// map32
		ans := make([]byte, 5)
		ans[0] = 0xdf
		binary.BigEndian.PutUint32(ans[1:5], uint32(l))
		return ans
	}
}

func getArrayFormat(v []interface{}) []byte {
	l := len(v)
	switch {
	case l < 16:
		// fixarray
		return []byte{(0x90 | byte(len(v)))}
	case l < int(math.Pow(2, 16)):
		// array16
		ans := make([]byte, 3)
		ans[0] = 0xdc
		binary.BigEndian.PutUint16(ans[1:3], uint16(l))
		return ans
	default:
		// array32
		ans := make([]byte, 5)
		ans[0] = 0xdd
		binary.BigEndian.PutUint32(ans[1:5], uint32(l))
		return ans
	}
}
