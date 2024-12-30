package client

import (
	"dynamicTcp/client"
	"dynamicTcp/protocol"
	"errors"
	"net"
	"time"
)

type DynamicTcpConn struct {
	dataBuf  []byte
	streamId int64
	client   *Client
}

func NewDynamicTcpConn(streamId int64, dataChan chan *protocol.DynamicTcpMessage, client *client.Client) *DynamicTcpConn {
	return &DynamicTcpConn{
		streamId: streamId,
		dataChan: dataChan,
		client:   client,
	}
}

// Read reads data from the connection.
// Read reads data from the connection.
func (c *DynamicTcpConn) Read(b []byte) (n int, err error) {
	data, ok := <-c.dataChan
	if !ok {
		return 0, errors.New("connection closed")
	}

	n = copy(b, data.StreamData)
	return n, nil
}

// Write writes data to the connection.
func (c *DynamicTcpConn) Write(b []byte) (n int, err error) {
	packet := protocol.MakeDynamicTcpBytes(c.streamId, b)
	return c.client.conn.Write(packet)
}

// Close closes the connection.
func (c *DynamicTcpConn) Close() error {
	return nil
}

// LocalAddr returns the local network address, if known.
func (c *DynamicTcpConn) LocalAddr() net.Addr {
	return c.client.conn.LocalAddr()
}

// RemoteAddr returns the remote network address, if known.
func (c *DynamicTcpConn) RemoteAddr() net.Addr {
	return c.client.conn.LocalAddr()
}

// SetDeadline sets the read and write deadlines associated with the connection.
func (c *DynamicTcpConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls and any currently-blocked Read call.
func (c *DynamicTcpConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls and any currently-blocked Write call.
func (c *DynamicTcpConn) SetWriteDeadline(t time.Time) error {
	return nil
}
