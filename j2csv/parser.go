package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

type Parser struct {
	headers    []string
	out        io.WriteCloser
	inp        io.ReadCloser
	decoder    *json.Decoder
	utsHeaders map[string]struct{}
}

func (p *Parser) writeRow(row map[string]any) {
	for index, header := range p.headers {
		value := row[header]

		switch v := value.(type) {
		case float64:
			if _, ok := p.utsHeaders[header]; ok {
				t := time.Unix(int64(v), 0)
				p.out.Write([]byte(t.String()))
			} else {
				p.out.Write([]byte(fmt.Sprintf("%v", v)))
			}
		default:
			p.out.Write([]byte(fmt.Sprintf("%v", v)))
		}

		if index < len(p.headers)-1 {
			p.out.Write([]byte(","))
		}

	}

	p.out.Write([]byte("\n"))
}

func (p *Parser) parseArray() {
	for p.decoder.More() {
		object := map[string]any{}
		if err := p.decoder.Decode(&object); err != nil {
			log.Fatalf("error while decoding object : %v", err)
		}

		p.writeRow(object)
	}
}

func (p *Parser) getHeaderAndFirstRow() ([]string, map[string]any, error) {
	if p.decoder.More() {
		object := map[string]any{}
		if err := p.decoder.Decode(&object); err != nil {
			log.Fatalf("error while decoding object : %v", err)
		}

		headers := []string{}

		for key := range object {
			headers = append(headers, key)
		}

		return headers, object, nil
	}

	return nil, nil, fmt.Errorf("empty array or object")
}

func (p *Parser) setHeadersAndWriteFirstRow() *Parser {
	if p.isArray() {
		headers, row, err := p.getHeaderAndFirstRow()
		if err != nil {
			log.Fatalf("error while getting headers : %v", err)
		}
		p.headers = headers
		p.writeRow(row)
		return p
	} else {
		log.Fatalf("Invalid Json")
	}
	return p
}

func (p *Parser) setUTS(uts string) *Parser {

	fields := strings.Split(uts, ",")

	if len(fields) <= 0 {
		return p
	}

	p.utsHeaders = make(map[string]struct{})

	for _, v := range fields {
		p.utsHeaders[v] = struct{}{}
	}

	return p
}

func (p *Parser) isArray() bool {
	token, err := p.decoder.Token()
	if err != nil {
		log.Fatalf("error while reading token : %v", err)
	}

	return token == json.Delim('[')
}

// func (p *Parser) isObject() bool {
// 	token, err := p.decoder.Token()
// 	if err != nil {
// 		log.Fatalf("error while reading start token : %v", err)
// 	}

// 	return token == "{"
// }
