package client

import (
	"dynamicTcp/protocol"
	"log"
	"net"
	"sync"
	"sync/atomic"
)

type Client struct {
	conn         net.Conn
	multipleConn map[int64]*DynamicTcpConn
	sessionID    int64
	mu           sync.Mutex
}

func NewClient() *Client {
	return &Client{
		multipleConn: make(map[int64]*DynamicTcpConn),
	}
}

func (c *Client) Start(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	c.conn = conn

	dataChan := protocol.ReadTcpStream(conn)
	go c.receiveData(dataChan)
}

func (c *Client) generateStreamId() int64 {
	return atomic.AddInt64(&c.sessionID, 1)
}

func (c *Client) receiveData(dataChan chan *protocol.DynamicTcpMessage) {
	for data := range dataChan {
		c.mu.Lock()
		if multipleConn, ok := c.multipleConn[data.StreamId]; ok {
			multipleConn.dataChan
			delete(c.multipleConn, data.StreamId)
		}
		c.mu.Unlock()
	}
}

func (c *Client) AddMultiConn() net.Conn {
	streamId := c.generateStreamId()

	// 여기서는 예시로 nil을 반환합니다.
	return nil
}

// Listener 인터페이스 구현
func (c *Client) OnMessage(msg *protocol.DynamicTcpMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if sessionChan, ok := c.multipleConn[msg.StreamId]; ok {
		sessionChan <- msg
	}
}

func (c *Client) OnError(err error) {
	log.Printf("Error: %v", err)
}
