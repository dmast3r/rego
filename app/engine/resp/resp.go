package resp

import (
	"bytes"
	"errors"
	"github.com/dmast3r/rego/app/config"
	"io"
	"net"
	"strconv"
)

type Parser struct {
	conn net.Conn
	buf  *bytes.Buffer
}

func NewRespParser(conn net.Conn) *Parser {
	return &Parser{
		conn: conn,
		buf:  bytes.NewBuffer([]byte{}),
	}
}

func (rp *Parser) DecodeRESP() (interface{}, error) {
	tempBuf := make([]byte, config.RESPBufferLength)
	n, err := rp.conn.Read(tempBuf)

	if err != nil {
		return nil, err
	}

	if n <= 0 {
		return nil, io.EOF
	}

	rp.buf.Write(tempBuf)
	return rp.processBuffer()
}

func (rp *Parser) processBuffer() (interface{}, error) {
	b, err := rp.buf.ReadByte()
	if err != nil {
		return "", err
	}

	switch b {
	case '+':
		return readSimpleString(rp)
	case '-':
		return readError(rp)
	case ':':
		return readInt64(rp)
	case '$':
		return readBulkString(rp)
	case '*':
		return readArray(rp)
	}

	return nil, errors.New("unexpected input found. The input must follow the REdis Serialisation Protocol")
}

func readDelimitedString(rp *Parser) (string, error) {
	s, err := rp.buf.ReadString('\r')
	if err != nil {
		return "", err
	}

	rp.buf.ReadByte() // to skip over the '\n' character
	return s[:len(s)-1], nil
}

func readSimpleString(rp *Parser) (string, error) {
	return readDelimitedString(rp)
}

func readError(rp *Parser) (string, error) {
	return readDelimitedString(rp)
}

func readInt64(rp *Parser) (int64, error) {
	s, err := readDelimitedString(rp)
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return v, nil
}

func readBulkString(rp *Parser) (string, error) {
	length, err := readInt64(rp)
	if err != nil {
		return "", nil
	}

	bulkStr := make([]byte, length)
	_, err = rp.buf.Read(bulkStr)
	if err != nil {
		return "", err
	}

	rp.buf.ReadByte()
	rp.buf.ReadByte()

	return string(bulkStr), nil
}

func readArray(rp *Parser) (interface{}, error) {
	length, err := readInt64(rp)
	if err != nil {
		return nil, err
	}

	elems := make([]interface{}, length)
	for i := range elems {
		elem, err := rp.processBuffer()
		if err != nil {
			return nil, err
		}
		elems[i] = elem
	}

	return elems, nil
}
