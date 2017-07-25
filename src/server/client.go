package server

import (
	"io"
	"sync"

	"github.com/gorilla/websocket"
)

// Client is a ReadWriteCloser wrapper around a websocket connection
type Client struct {
	conn       *websocket.Conn
	reader     io.Reader
	writer     io.WriteCloser
	writeMutex sync.Mutex
}

// Read implements ReadWriteCloser interface
func (c *Client) Read(p []byte) (n int, err error) {
	if c.reader == nil {
		_, c.reader, err = c.conn.NextReader()
		if err != nil {
			n = 0
			return
		}
	}
	for n = 0; n < len(p); {
		var bytes int
		bytes, err = c.reader.Read(p[n:])
		n += bytes
		if err == io.EOF {
			c.reader = nil
			break
		}
		if err != nil {
			break
		}
	}
	return
}

// Write implements ReadWriteCloser interface
func (c *Client) Write(p []byte) (n int, err error) {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	if c.writer == nil {
		c.writer, err = c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			n = 0
			return
		}
	}
	for n = 0; n < len(p); {
		var bytes int
		bytes, err = c.writer.Write(p)
		n += bytes
		if err != nil {
			break
		}
	}
	if err != nil || n == len(p) {
		err = c.Close()
	}
	return
}

// WriteJSON is a wrapper for Websocket.WriteJSON to ensure no concurrent writes
func (c *Client) WriteJSON(v interface{}) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	return c.conn.WriteJSON(v)
}

// Close implements ReadWriteCloser interface
func (c *Client) Close() (err error) {
	if c.writer != nil {
		err = c.writer.Close()
		c.writer = nil
	}
	return
}
