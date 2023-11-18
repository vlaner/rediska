package resp

import (
	"strconv"
)

const (
	STRING      = '+'
	ERROR       = '-'
	BULK_STRING = '$'
	INTEGER     = ':'
	ARRAY       = '*'
)

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

func (v Value) Marshal() []byte {
	switch v.Typ {
	case "string":
		return v.marshalString()
	case "bulk":
		return v.marshalBulkString()
	case "num":
		return v.marshalInt()
	case "array":
		return v.marshalArray()
	case "error":
		return v.marshalError()
	case "null":
		return v.marshalNull()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var b []byte
	b = append(b, STRING)
	b = append(b, v.Str...)
	b = append(b, '\r', '\n')
	return b
}

func (v Value) marshalError() []byte {
	var b []byte
	b = append(b, ERROR)
	b = append(b, v.Str...)
	b = append(b, '\r', '\n')
	return b
}

func (v Value) marshalBulkString() []byte {
	var b []byte
	b = append(b, BULK_STRING)
	b = append(b, strconv.Itoa(len(v.Bulk))...)
	b = append(b, '\r', '\n')
	b = append(b, []byte(v.Bulk)...)
	b = append(b, '\r', '\n')
	return b
}

func (v Value) marshalInt() []byte {
	var b []byte
	b = append(b, INTEGER)
	b = append(b, strconv.Itoa(v.Num)...)
	b = append(b, '\r', '\n')
	return b
}

func (v Value) marshalArray() []byte {
	var b []byte
	b = append(b, ARRAY)
	b = append(b, strconv.Itoa(len(v.Array))...)
	b = append(b, '\r', '\n')

	for _, val := range v.Array {
		b = append(b, val.Marshal()...)
	}

	return b
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}
