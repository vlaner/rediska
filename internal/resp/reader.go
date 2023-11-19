package resp

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strconv"
)

var (
	ErrInvalidBulkStringSize = errors.New("invalid bulk string size")
	ErrInvalidArraySize      = errors.New("invalid array size")
	ErrEmptyInput            = errors.New("invalid input")
)

type Reader struct {
	reader *textproto.Reader
}

func NewReader(r io.Reader) *Reader {

	rd := bufio.NewReader(r)
	return &Reader{
		reader: textproto.NewReader(rd),
	}
}

func (r *Reader) readLine() ([]byte, error) {
	return r.reader.ReadLineBytes()
}

func (r *Reader) ParseInput() (Value, error) {
	typ, err := r.reader.R.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch typ {
	case BULK_STRING:
		return r.parseBulkString()
	case INTEGER:
		return r.readInteger()
	case ARRAY:
		return r.parseArray()
	default:
		return Value{}, fmt.Errorf("datatype not supported: %v", string(typ))
	}
}

func (r *Reader) getLength() (int, error) {
	line, err := r.readLine()
	if err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(string(line))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (r *Reader) parseBulkString() (Value, error) {
	v := Value{Typ: "bulk"}

	strLen, err := r.getLength()
	if err != nil {
		return v, err
	}
	if strLen < 0 {
		return v, ErrInvalidBulkStringSize
	}

	buf := make([]byte, strLen)
	r.reader.R.Read(buf)
	v.Bulk = string(buf)

	// read trailing CRLF after string
	r.readLine()

	return v, nil
}

func (r *Reader) parseArray() (Value, error) {
	v := Value{Typ: "array"}

	arrayLen, err := r.getLength()
	if err != nil {
		return v, err
	}
	if arrayLen < 0 {
		return v, ErrInvalidArraySize
	}
	for i := 0; i < arrayLen; i++ {
		elem, err := r.ParseInput()
		if err != nil {
			return v, err
		}
		v.Array = append(v.Array, elem)
	}

	return v, nil
}

func (r *Reader) readInteger() (Value, error) {
	v := Value{Typ: "num"}

	line, err := r.readLine()
	if err != nil {
		return v, err
	}
	n, err := strconv.Atoi(string(line))
	if err != nil {
		return v, err
	}
	v.Num = n
	return v, nil
}
