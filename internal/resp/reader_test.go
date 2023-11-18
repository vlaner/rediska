package resp_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/vlaner/rediska/internal/resp"
)

func TestProtocolParseSimpleBulkString(t *testing.T) {
	input := []byte("$4\r\nTEST\r\n")
	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)
	output, err := r.ParseInput()
	if err != nil {
		t.Fatalf("Expected no error got %v", err)
	}
	result := output.Marshal()
	if !bytes.Equal(result, input) {
		t.Errorf("EXPECTED OUTPUT TO BE '%v' BUT GOT %v", input, result)
	}
}

func TestProtocolParseArrayOfBulkStrings(t *testing.T) {
	input := []byte("*3\r\n$4\r\nTEST\r\n$5\r\nTEST1\r\n$8\r\nEIGHTLEN\r\n")
	buf := new(bytes.Buffer)
	buf.Write(input)

	r := resp.NewReader(buf)

	output, err := r.ParseInput()
	if err != nil {
		t.Fatalf("Expected no error got %v", err)
	}

	result := output.Marshal()
	if !bytes.Equal(result, input) {
		t.Errorf("EXPECTED OUTPUT TO BE '%v' BUT GOT %v", input, result)
	}
}

func TestProtocolParseErrorNoPrefix(t *testing.T) {
	input := []byte("4\r\nTEST\r\n")
	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)

	_, err := r.ParseInput()
	if err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestProtocolParseErrorWrongBulkStringLength(t *testing.T) {
	input := []byte("$-1\r\nTEST\r\n")

	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)

	_, err := r.ParseInput()
	if !errors.Is(err, resp.ErrInvalidBulkStringSize) {
		t.Fatalf("Expected %v error got %v", resp.ErrInvalidBulkStringSize, err)
	}
}

func TestProtocolParseErrorWrongBulkStringSize(t *testing.T) {
	input := []byte("$5\r\naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\r\n")

	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)

	_, err := r.ParseInput()
	if err != nil {
		t.Fatalf("Unexpected error got %v", err)
	}
}

func TestProtocolParseErrorWrongBulkStringSizeNotNumber(t *testing.T) {
	input := []byte("$x\r\na\r\n")

	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)

	_, err := r.ParseInput()
	if err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestProtocolParseErrorWrongArraySizeNotNumber(t *testing.T) {
	input := []byte("*x\r\n")

	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)

	_, err := r.ParseInput()
	if err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestProtocolParseErrorWrongArraySize(t *testing.T) {
	input := []byte("*-1\r\n$4\r\nTEST\r\n$5\r\nTEST1\r\n$8\r\nEIGHTLEN\r\n")

	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)

	_, err := r.ParseInput()
	if !errors.Is(err, resp.ErrInvalidArraySize) {
		t.Fatalf("Expected error got %v", err)
	}
}

func TestProtocolParseEmptyInput(t *testing.T) {
	var input []byte
	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)

	_, err := r.ParseInput()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Expected %v error got %v", io.EOF, err)
	}
}

func TestProtocolParseArrayOfBulkStringsWithWrongSize(t *testing.T) {
	input := []byte("*3\r\n$7\r\nTEST\r\n$5\r\nTEST1\r\n$8\r\nEIGHTLEN\r\n")

	buf := new(bytes.Buffer)
	buf.Write(input)
	r := resp.NewReader(buf)

	_, err := r.ParseInput()
	if err == nil {
		t.Fatalf("Expected error got %v", err)
	}
}
