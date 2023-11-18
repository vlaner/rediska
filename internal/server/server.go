package server

import (
	"context"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/vlaner/rediska/internal/commands"
	"github.com/vlaner/rediska/internal/resp"
)

type Server struct {
	listenAddr string
	l          net.Listener
	quit       chan struct{}
	wg         *sync.WaitGroup
}

func New(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		l:          nil,
		quit:       make(chan struct{}),
		wg:         &sync.WaitGroup{},
	}
}

func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.l = l
	s.wg.Add(1)
	go s.acceptLoop(ctx)

	return nil
}

func (s *Server) Stop() {
	close(s.quit)
	s.l.Close()
	s.wg.Wait()
}

func (s Server) acceptLoop(ctx context.Context) {
	defer s.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.quit:
			return
		default:
			conn, err := s.l.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				case <-ctx.Done():
					return
				default:
					log.Printf("error accepting connection: %v\n", err)
				}
			}
			s.wg.Add(1)
			go func() {
				s.handleConnection(conn)
				defer s.wg.Done()
			}()
		}
	}
}

func (s Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	r := resp.NewReader(conn)
	w := resp.NewWriter(conn)

	for {
		v, err := r.ParseInput()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
		}

		if err != nil {
			w.Write(resp.Value{Typ: "error", Str: err.Error()})
			continue
		}

		cmd := strings.ToUpper(v.Array[0].Bulk)
		args := v.Array[1:]
		handler, ok := commands.Handlers[cmd]
		if !ok {
			w.Write(resp.Value{Typ: "error", Str: "command does not exist"})
			continue
		}

		result := handler(args)

		w.Write(result)
	}
}
