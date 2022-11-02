package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type client struct {
	conn      *websocket.Conn
	logStream io.ReadCloser
}

func (c *client) readPump(ctx context.Context, stopCh chan struct{}) {
	defer func() {
		c.conn.Close()
		stopCh <- struct{}{}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, _, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
					log.Printf("error: %v", err)
				}
				return
			}
		}
	}
}

func (c *client) writePump(ctx context.Context, stopCh chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		stopCh <- struct{}{}
	}()

	reader := bufio.NewReader(c.logStream)
	for {
		select {
		case <-ctx.Done():
		case <-stopCh:
			return
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		default:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			line, err := readLongLine(reader)
			if err != nil {
				log.Println("writePump: cannot read line", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, line); err != nil {
				log.Println("writePump: cannot write message", err)
				continue
			}
		}
	}
}

func readLongLine(r *bufio.Reader) (line []byte, err error) {
	var buffer []byte
	var isPrefix bool

	for {
		buffer, isPrefix, err = r.ReadLine()
		line = append(line, buffer...)
		if err != nil {
			break
		}

		if !isPrefix {
			break
		}
	}

	return line, err
}
