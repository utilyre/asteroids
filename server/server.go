package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
)

type Server struct {
	Addr   string
	LogDir string

	handlersMap map[string]Handler
}

func (srv *Server) ListenAndServe() error {
	addr := "localhost:80"
	if len(srv.Addr) > 0 {
		addr = srv.Addr
	}
	if srv.handlersMap == nil {
		srv.handlersMap = map[string]Handler{}
	}

	logger := slog.With("address", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Info("started listening")

	// serve
	for {
		logger.Debug("waiting to accept a new connection")
		conn, err := ln.Accept()
		if err != nil {
			logger.Error("failed to establish connection", "error", err)
			continue
		}
		logger.Info("connection accepted", "remote", conn.RemoteAddr())

		monitorBuffer := &bytes.Buffer{}
		connReader := io.TeeReader(conn, monitorBuffer)

		logger.Debug("monitoring connection")
		go srv.monitorConn(monitorBuffer, conn.RemoteAddr())

		logger.Debug("handling connection")
		go func() {
			defer func() {
				if err := conn.Close(); err != nil {
					logger.Error("failed to close connection", "error", err)
				}
			}()

			srv.handleConn(connReader, conn.RemoteAddr())
		}()
	}
}

func (srv *Server) handleConn(r io.Reader, remote net.Addr) {
	logger := slog.With("remote", remote)

	for {
		logger.Debug("reading message from network")

		msg, err := ReadMessage(r)
		if errors.Is(err, io.EOF) {
			logger.Info("connection closed")
			break
		}
		if err != nil {
			logger.Error("failed to read message from connection", "error", err)
			continue
		}
		logger := logger.With(slog.Group("message", "version", msg.Version, "scope", msg.Scope))
		logger.Debug("message received")

		if err := srv.dispatchMessage(msg); err != nil {
			logger.Warn("failed to dispatch message")
		}
	}
}

type Handler func(ctx context.Context, body []byte) error

func (srv *Server) Handle(scope string, handler Handler) {
	if strings.ContainsRune(scope, ' ') {
		panic("scope cannot contain spaces")
	}

	if srv.handlersMap == nil {
		srv.handlersMap = map[string]Handler{}
	}

	if _, exists := srv.handlersMap[scope]; exists {
		panic("scope already handled")
	}

	srv.handlersMap[scope] = handler
}

func (srv *Server) dispatchMessage(msg *Message) error {
	if msg.Version != 1 {
		return fmt.Errorf("unsupported message version: %d", msg.Version)
	}

	handler, exists := srv.handlersMap[msg.Scope]
	if !exists {
		return fmt.Errorf("no handler specified for scope: %s", msg.Scope)
	}

	if err := handler(context.TODO(), msg.Body); err != nil {
		return fmt.Errorf("handle %s: %w", msg.Scope, err)
	}

	return nil

	/* switch msg.Scope {
	case "player.move_forward":
		slog.Debug("player moved forward")
	case "player.move_backward":
		slog.Debug("player moved backward")
	case "player.rotate_left":
		slog.Debug("player rotated left")
	case "player.rotate_right":
		slog.Debug("player rotated right")
	case "player.shoot":
		slog.Debug("player shot")
	} */
}

func (srv *Server) monitorConn(r io.Reader, remote net.Addr) {
	logger := slog.With("remote", remote)

	logDir := "logs"
	if len(srv.LogDir) > 0 {
		logDir = srv.LogDir
	}

	if err := os.MkdirAll(logDir, 0o755); err != nil {
		logger.Error("failed to make log directory", "path", logDir, "error", err)
		return
	}

	name := fmt.Sprintf("logs/traffic_%s.log", remote)
	logger = logger.With("filename", name)

	f, err := os.Create(name)
	if err != nil {
		logger.Error("failed to create log file", "error", err)
		return
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	logger.Info("copying traffic to log file")
	_, err = io.Copy(f, r)
	if err != nil {
		logger.Error("failed to copy traffic to log file", "error", err)
	}
}

func readUInt64(r io.Reader) (uint64, error) {
	var value uint64
	if err := binary.Read(r, binary.BigEndian, &value); err != nil {
		return 0, err
	}

	return value, nil
}

func writeUInt64(w io.Writer, value uint64) error {
	if err := binary.Write(w, binary.BigEndian, value); err != nil {
		return err
	}

	return nil
}

func readBytes(r io.Reader) ([]byte, error) {
	size, err := readUInt64(r)
	if err != nil {
		return nil, err
	}
	if size > 1024 {
		return nil, errors.New("receiving over-sized bytes")
	}

	buf := make([]byte, size)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	if uint64(n) < size {
		panic("read bytes and pre-known size mismatch")
	}

	return buf, nil
}

func writeBytes(w io.Writer, p []byte) error {
	size := uint64(len(p))
	if err := writeUInt64(w, size); err != nil {
		return err
	}

	n, err := w.Write(p)
	if err != nil {
		return err
	}
	if uint64(n) < size {
		panic("written bytes and the size of bytes mismatch")
	}

	return nil
}
