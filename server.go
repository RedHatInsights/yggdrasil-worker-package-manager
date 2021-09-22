package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"git.sr.ht/~spc/go-log"
	"github.com/google/uuid"
	"github.com/redhatinsights/yggdrasil/protocol"
	"google.golang.org/grpc"
)

type Message struct {
	Command string `json:"command"`
	Name    string `json:"name"`
}

type Output struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Code   int    `json:"code"`
}

type Server struct {
	protocol.UnimplementedWorkerServer

	dispatchAddr  string
	allowPatterns []*regexp.Regexp
	pm            PackageManager
}

func (s *Server) Send(ctx context.Context, d *protocol.Data) (*protocol.Receipt, error) {
	log.Debugf("received message: %v", d.MessageId)
	log.Tracef("%v", d)
	go func() {
		var m Message
		if err := json.Unmarshal(d.GetContent(), &m); err != nil {
			log.Errorf("cannot unmarshal data: %v", err)
			return
		}
		log.Debugf("received command: %v", m)

		if !s.packageAllowed(m.Name) {
			log.Errorf("cannot install %v: does not match an allow pattern", m.Name)
			return
		}

		var stdout, stderr []byte
		var code int
		var err error
		switch m.Command {
		case "install":
			stdout, stderr, code, err = s.pm.Install(m.Name)
			if err != nil {
				log.Errorf("cannot install package: %v", err)
				return
			}
			log.Infof("installed package: %v", m.Name)
		case "remove":
			stdout, stderr, code, err = s.pm.Uninstall(m.Name)
			if err != nil {
				log.Errorf("cannot remove package: %v", err)
				return
			}
			log.Infof("removed package: %v", m.Name)
		default:
			log.Errorf("cannot perform command: %v", m.Command)
			return
		}

		output := Output{
			Stdout: string(stdout),
			Stderr: string(stderr),
			Code:   code,
		}
		log.Trace(output.Stdout)
		if output.Code != 0 {
			log.Trace(output.Stderr)
		}

		content, err := json.Marshal(output)
		if err != nil {
			log.Errorf("cannot marshal data: %v", err)
			return
		}

		data := &protocol.Data{
			MessageId:  uuid.New().String(),
			Metadata:   map[string]string{},
			Content:    content,
			ResponseTo: d.GetMessageId(),
			Directive:  "",
		}
		if err := s.returnData(data); err != nil {
			log.Errorf("cannot return data: %v", err)
			return
		}
		log.Debugf("published message: %v", data.MessageId)
		log.Tracef("%v", data)
	}()

	return &protocol.Receipt{}, nil
}

func (s *Server) returnData(d *protocol.Data) error {
	conn, err := grpc.Dial(s.dispatchAddr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("cannot dial dispatcher: %w", err)
	}
	defer conn.Close()

	c := protocol.NewDispatcherClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if _, err := c.Send(ctx, d); err != nil {
		return fmt.Errorf("cannot send data: %w", err)
	}

	return nil
}

func (s *Server) packageAllowed(name string) bool {
	for _, r := range s.allowPatterns {
		if r.Match([]byte(name)) {
			return true
		}
	}
	return false
}