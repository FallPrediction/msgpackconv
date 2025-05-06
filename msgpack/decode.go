package msgpack

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"reflect"
	"slices"
)

var LastByte = map[string]byte{
	"positiveFixint": 0x7f,
	"negativeFixint": 0xff,
	"fixstr":         0xbf,
	"fixarray":       0x9f,
	"fixmap":         0x8f,
}

func ToJSON(msgpackconv []byte) []byte {
	if len(msgpackconv) == 0 {
		return []byte{}
	}
	obj, _, err := decode(msgpackconv)
	if err == ErrInvalidMsgPack {
		return []byte{}
	}
	ans, err := json.Marshal(obj)
	if err == ErrInvalidMsgPack {
		return []byte{}
	}
	return ans
}

func decode(msgpackconv []byte) (interface{}, int, error) {
	var obj interface{}
	v := reflect.ValueOf(&obj).Elem()
	// 已讀取的 byte index
	idxOfEnd := 0
	if msgpackconv[0] >= FirstByte["fixstr"] && msgpackconv[0] <= LastByte["fixstr"] {
		// fixstr
		l := getLength([]byte{msgpackconv[0] ^ FirstByte["fixstr"]})
		if len(msgpackconv) < l+1 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(string(msgpackconv[1 : 1+l])))
		idxOfEnd = l + 1
	} else if msgpackconv[0] == FirstByte["str8"] {
		// str8
		if len(msgpackconv) < 2 {
			return nil, 0, ErrInvalidMsgPack
		}
		l := getLength(msgpackconv[1:2])
		if len(msgpackconv) < l+2 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(string(msgpackconv[2 : 2+l])))
		idxOfEnd = l + 2
	} else if msgpackconv[0] == FirstByte["str16"] {
		// str16
		if len(msgpackconv) < 3 {
			return nil, 0, ErrInvalidMsgPack
		}
		l := getLength(msgpackconv[1:3])
		if len(msgpackconv) < l+3 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(string(msgpackconv[3 : 3+l])))
		idxOfEnd = l + 3
	} else if msgpackconv[0] == FirstByte["str32"] {
		// str32
		if len(msgpackconv) < 5 {
			return nil, 0, ErrInvalidMsgPack
		}
		l := getLength(msgpackconv[1:5])
		if len(msgpackconv) < l+5 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(string(msgpackconv[5 : 5+l])))
		idxOfEnd = l + 5
	} else if msgpackconv[0] >= FirstByte["positiveFixint"] && msgpackconv[0] <= LastByte["positiveFixint"] {
		// positive fixint
		v.Set(reflect.ValueOf(bytesToFloat64([]byte{msgpackconv[0]}, true)))
		idxOfEnd = 1
	} else if msgpackconv[0] == FirstByte["uint8"] {
		// uint8
		if len(msgpackconv) < 2 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bytesToFloat64(msgpackconv[1:2], true)))
		idxOfEnd = 2
	} else if msgpackconv[0] == FirstByte["uint16"] {
		// uint16
		if len(msgpackconv) < 3 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bytesToFloat64(msgpackconv[1:3], true)))
		idxOfEnd = 3
	} else if msgpackconv[0] == FirstByte["uint32"] {
		// uint32
		if len(msgpackconv) < 5 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bytesToFloat64(msgpackconv[1:5], true)))
		idxOfEnd = 5
	} else if msgpackconv[0] == FirstByte["uint64"] {
		// uint64
		if len(msgpackconv) < 9 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bytesToFloat64(msgpackconv[1:9], true)))
		idxOfEnd = 9
	} else if msgpackconv[0] >= FirstByte["negativeFixint"] {
		// negative fixint
		v.Set(reflect.ValueOf(bytesToFloat64([]byte{msgpackconv[0]}, false)))
		idxOfEnd = 1
	} else if msgpackconv[0] == FirstByte["int8"] {
		// int8
		if len(msgpackconv) < 2 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bytesToFloat64(msgpackconv[1:2], false)))
		idxOfEnd = 2
	} else if msgpackconv[0] == FirstByte["int16"] {
		// int16
		if len(msgpackconv) < 3 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bytesToFloat64(msgpackconv[1:3], false)))
		idxOfEnd = 3
	} else if msgpackconv[0] == FirstByte["int32"] {
		// int32
		if len(msgpackconv) < 5 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bytesToFloat64(msgpackconv[1:5], false)))
		idxOfEnd = 5
	} else if msgpackconv[0] == FirstByte["int64"] {
		// int64
		if len(msgpackconv) < 9 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bytesToFloat64(msgpackconv[1:9], false)))
		idxOfEnd = 9
	} else if msgpackconv[0] == FirstByte["float32"] {
		// float32
		if len(msgpackconv) < 5 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bitsToFloat32(msgpackconv[1:5])))
		idxOfEnd = 5
	} else if msgpackconv[0] == FirstByte["float64"] {
		// float64
		if len(msgpackconv) < 9 {
			return nil, 0, ErrInvalidMsgPack
		}
		v.Set(reflect.ValueOf(bitsToFloat64(msgpackconv[1:9])))
		idxOfEnd = 9
	} else if msgpackconv[0] == FirstByte["nil"] {
		// nil
		// do nothing
		idxOfEnd = 1
	} else if msgpackconv[0] == FirstByte["false"] {
		// false
		v.Set(reflect.ValueOf(false))
		idxOfEnd = 1
	} else if msgpackconv[0] == FirstByte["true"] {
		// true
		v.Set(reflect.ValueOf(true))
		idxOfEnd = 1
	} else if msgpackconv[0] >= FirstByte["fixarray"] && msgpackconv[0] <= LastByte["fixarray"] {
		// fixarray
		l := getLength([]byte{msgpackconv[0] ^ FirstByte["fixarray"]})
		s := make([]interface{}, l)
		var err error
		for j := range l {
			s[j], _, err = decode(msgpackconv[1+j : 2+j])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
		}
		v.Set(reflect.ValueOf(s))
		idxOfEnd = l + 1
	} else if msgpackconv[0] == FirstByte["array16"] {
		// array16
		l := getLength(msgpackconv[1:3])
		s := make([]interface{}, l)
		var err error
		for j := range l {
			s[j], _, err = decode(msgpackconv[3+j : 4+j])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
		}
		v.Set(reflect.ValueOf(s))
		idxOfEnd = l + 3
	} else if msgpackconv[0] == FirstByte["array32"] {
		// array32
		l := getLength(msgpackconv[1:5])
		s := make([]interface{}, l)
		var err error
		for j := range l {
			s[j], _, err = decode(msgpackconv[5+j : 6+j])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
		}
		v.Set(reflect.ValueOf(s))
		idxOfEnd = l + 5
	} else if msgpackconv[0] >= FirstByte["fixmap"] && msgpackconv[0] <= LastByte["fixmap"] {
		// fixmap
		m := make(map[string]interface{})
		l := getLength([]byte{msgpackconv[0] ^ FirstByte["fixmap"]})
		j := 0
		for len(m) < l {
			key, tmp, err := decode(msgpackconv[j+1:])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
			j += tmp
			value, tmp, err := decode(msgpackconv[j+1:])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
			j += tmp
			m[key.(string)] = value
		}
		v.Set(reflect.ValueOf(m))
		idxOfEnd = j
	} else if msgpackconv[0] == FirstByte["map16"] {
		// map16
		m := make(map[string]interface{})
		l := getLength(msgpackconv[1:3])
		j := 2
		for len(m) < l {
			key, tmp, err := decode(msgpackconv[j+1:])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
			j += tmp
			value, tmp, err := decode(msgpackconv[j+1:])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
			j += tmp
			m[key.(string)] = value
		}
		v.Set(reflect.ValueOf(m))
		idxOfEnd = j
	} else if msgpackconv[0] == FirstByte["map32"] {
		// map32
		m := make(map[string]interface{})
		l := getLength(msgpackconv[1:5])
		j := 4
		for len(m) < l {
			key, tmp, err := decode(msgpackconv[j+1:])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
			j += tmp
			value, tmp, err := decode(msgpackconv[j+1:])
			if err != nil {
				return nil, 0, ErrInvalidMsgPack
			}
			j += tmp
			m[key.(string)] = value
		}
		v.Set(reflect.ValueOf(m))
		idxOfEnd = j
	}
	return obj, idxOfEnd, nil
}

func getLength(bytes []byte) int {
	return int(bytesToUint64(bytes, true))
}

func bytesToFloat64(bytes []byte, positive bool) float64 {
	return float64(int64(bytesToUint64(bytes, positive)))
}

func bytesToUint64(bytes []byte, positive bool) uint64 {
	var copyBytes []byte
	if positive {
		copyBytes = slices.Repeat([]byte{FirstByte["positiveFixint"]}, 8)
	} else {
		copyBytes = slices.Repeat([]byte{LastByte["negativeFixint"]}, 8)
	}
	copy(copyBytes[8-len(bytes):], bytes)
	return binary.BigEndian.Uint64(copyBytes)
}

func bitsToFloat64(bytes []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(bytes))
}

func bitsToFloat32(bytes []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(bytes))
}
