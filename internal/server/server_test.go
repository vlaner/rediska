package server_test

import (
	"bytes"
	"context"
	"net"
	"testing"

	"github.com/vlaner/rediska/internal/server"
)

func init() {
	s := server.New(":3000")
	ctx := context.Background()
	s.Start(ctx)
}

func TestPing(t *testing.T) {
	input := []byte("*2\r\n$4\r\nping\r\n$4\r\nTEST\r\n")

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		t.Fatalf("Dial failed: %s", err.Error())
	}
	defer conn.Close()

	_, err = conn.Write(input)
	if err != nil {
		t.Fatalf("Write to server failed: %s", err.Error())
	}
	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		t.Fatalf("Write to server failed: %s", err.Error())
	}
	expected := []byte("+TEST\r\n")
	if !bytes.Equal(reply[:n], expected) {
		t.Errorf("expected %v, got %v", expected, reply[:n])
	}
}

func TestGetNil(t *testing.T) {

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		t.Fatalf("Dial failed: %s", err.Error())
	}
	defer conn.Close()
	input := []byte("*2\r\n$3\r\nget\r\n$3\r\nkey\r\n")
	_, err = conn.Write(input)
	if err != nil {
		t.Fatalf("Write to server failed: %s", err.Error())
	}
	reply := make([]byte, 1024)
	n, err := conn.Read(reply)
	if err != nil {
		t.Fatalf("Write to server failed: %s", err.Error())
	}
	expected := []byte("$-1\r\n")
	if !bytes.Equal(reply[:n], expected) {
		t.Errorf("expected %v, got %v", expected, reply[:n])
	}
}

func TestSetAndGet(t *testing.T) {
	set := []byte("*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		t.Fatalf("Dial failed: %s", err.Error())
	}
	defer conn.Close()

	_, err = conn.Write(set)
	if err != nil {
		t.Fatalf("Write to server failed: %s", err.Error())
	}
	replySet := make([]byte, 1024)
	n, err := conn.Read(replySet)
	if err != nil {
		t.Fatalf("Write to server failed: %s", err.Error())
	}

	expectedSet := []byte("+OK\r\n")
	if !bytes.Equal(replySet[:n], expectedSet) {
		t.Errorf("expected %v, got %v", expectedSet, replySet[:n])
	}

	get := []byte("*2\r\n$3\r\nget\r\n$3\r\nkey\r\n")
	_, err = conn.Write(get)
	if err != nil {
		t.Fatalf("Write to server failed: %s", err.Error())
	}
	replyGet := make([]byte, 1024)
	n, err = conn.Read(replyGet)
	if err != nil {
		t.Fatalf("Write to server failed: %s", err.Error())
	}

	expectedGet := []byte("$5\r\nvalue\r\n")
	if !bytes.Equal(replyGet[:n], expectedGet) {
		t.Errorf("expected %v, got %v", expectedGet, replyGet[:n])
	}
}
