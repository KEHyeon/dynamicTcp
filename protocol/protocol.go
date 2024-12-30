package protocol

import (
	"encoding/binary"
	"io"
	"log"
)

// 헤더 상수 정의
const (
	HeaderSize     = 12
	StreamIdSize   = 8
	DataLengthSize = 4
)

// DynamicTcpMessage 구조체 정의
type DynamicTcpMessage struct {
	StreamId   int64
	StreamData []byte
	Err        error
}

func MakeDynamicTcpMessage(streamId int64, streamData []byte, err error) *DynamicTcpMessage {
	return &DynamicTcpMessage{
		StreamId:   streamId,
		StreamData: streamData,
		Err:        err,
	}
}
func MakeDynamicTcpBytes(streamId int64, streamData []byte) []byte {
	header := make([]byte, HeaderSize)
	binary.BigEndian.PutUint64(header[:StreamIdSize], uint64(streamId))
	binary.BigEndian.PutUint32(header[StreamIdSize:], uint32(len(streamData)))
	return append(header, streamData...)
}

// readTcpStream 함수 정의
func ReadTcpStream(reader io.Reader) chan *DynamicTcpMessage {
	dataChan := make(chan *DynamicTcpMessage)
	go func() {
		defer close(dataChan)
		for {
			// 헤더 읽기 (streamId: 8바이트, 데이터 길이: 4바이트)
			header := make([]byte, HeaderSize)
			_, err := io.ReadFull(reader, header)
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading header from connection: %v", err)
				}
				dataChan <- MakeDynamicTcpMessage(0, nil, err)
				return
			}

			// streamId와 데이터 길이 추출
			streamId := int64(binary.BigEndian.Uint64(header[:StreamIdSize]))
			dataLength := int(binary.BigEndian.Uint32(header[StreamIdSize:HeaderSize]))

			// 데이터 읽기
			data := make([]byte, dataLength)
			_, err = io.ReadFull(reader, data)
			if err != nil {
				if err != io.EOF {
					log.Printf("Error reading data from connection: %v", err)
				}
				dataChan <- MakeDynamicTcpMessage(streamId, nil, err)
				return
			}

			// 읽은 데이터를 채널에 전달
			dataChan <- MakeDynamicTcpMessage(streamId, data, nil)
		}
	}()
	return dataChan
}
