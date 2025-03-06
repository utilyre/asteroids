package main

import "io"

type Message struct {
	Version uint64
	Scope   string
	Body    []byte
}

func ReadMessage(r io.Reader) (*Message, error) {
	version, err := readUInt64(r)
	if err != nil {
		return nil, err
	}

	scope, err := readBytes(r)
	if err != nil {
		return nil, err
	}

	body, err := readBytes(r)
	if err != nil {
		return nil, err
	}

	return &Message{
		Version: version,
		Scope:   string(scope),
		Body:    body,
	}, nil
}
