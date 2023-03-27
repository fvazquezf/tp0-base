package common

import (
	"fmt"
	"net"
	"io"
	log "github.com/sirupsen/logrus"
)

// wrapper for socket connection that handles send/recv
type Socket struct {
	conn   net.Conn
}


// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func NewConnectedSocket(address string) (*Socket, error) {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
    }
    return &Socket{conn}, nil
}

func (s *Socket) sendSome(data []byte, sz uint32) (int, error) {
    sent := 0
    for sent < int(sz) {
        n, err := s.conn.Write(data[sent:sz])
        if err != nil {
            return 0, err
        }
        sent += n
    }
    return sent, nil
}

func (s *Socket) SendAll(data []byte) (error) {
    sent := 0
	sz := len(data)
    for sent < sz {
        n, err := s.sendSome(data[sent:sz], uint32(sz-sent))
        if err != nil {
            return err
        }
        sent += n
    }
    return nil
}

func (s *Socket) recvSome(data []byte, sz uint32) (int, error) {
    received := 0
    for received < int(sz) {
        n, err := s.conn.Read(data[received:sz])
        if err == io.EOF {
            received += n
            break
        }
        if err != nil {
            return 0, err
        }
        received += n
    }
    return received, nil
}

func (s *Socket) RecvAll(data []byte) error {
    received := 0
    sz := len(data)
    log.Infof("Received data: %d bytes", sz)
    for received < sz {
        n, err := s.recvSome(data[received:], uint32(sz-received))
        if err != nil {
            return err
        }
        received += n
    }
    return nil
}

func (s *Socket) Close() error {
    return s.conn.Close()
}