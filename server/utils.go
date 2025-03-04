package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log/slog"
)

func recvString(r io.Reader) (string, error) {
	var size uint64
	if err := binary.Read(r, binary.BigEndian, &size); err != nil {
		return "", err
	}
	if size > 64 {
		slog.Debug("size", "size", size)
		return "", errors.New("receiving an over-sized string")
	}

	slog.Info("read a string size", "size", size)

	buf := make([]byte, size)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	if uint64(n) < size {
		panic("read bytes and pre-known size mismatch")
	}

	return string(buf), nil
}

func sendString(w io.Writer, s string) error {
	buf := []byte(s)
	var size uint64 = uint64(len(buf))

	if err := binary.Write(w, binary.BigEndian, size); err != nil {
		return err
	}

	n, err := w.Write(buf)
	if err != nil {
		return err
	}
	if uint64(n) < size {
		panic("written bytes and the size of string mismatch")
	}

	return nil
}
