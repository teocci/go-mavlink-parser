// Package ws_net
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package ws_net

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsPort   = 7474
	serverIP = "192.168.100.92"
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}

	localIP = net.ParseIP(serverIP)

	localAddress  = fmt.Sprintf("localhost:%d", wsPort)
	remoteAddress = fmt.Sprintf("%s:%d", serverIP, wsPort)
)

// Client is a middleman between the ws-net connection and the hub.
type Client struct {
	// The ws-net connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send      chan []byte
	Done      chan struct{}
	Interrupt chan os.Signal
}

// onMessage reads messages from the ws-net connection
//
// The application runs onMessage in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) onMessage() {
	defer func() {
		_ = c.Conn.Close()
	}()

	fmt.Println("onMessage listening")
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}

// onEvent writes messages from the hub to the ws-net connection.
//
// A goroutine running onEvent is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) onEvent() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	fmt.Println("onEvent listening")
	for {
		select {
		case <-c.Done:
			return
		case message, ok := <-c.Send:
			data := string(message[:])
			log.Println("onEvent-> message")
			log.Printf("%#v", message)
			log.Printf("%#v", data)
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Println("onEvent-> !ok")
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current ws-net message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			log.Println("onEvent-> ticker.C")
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-c.Interrupt:
			log.Println("onEvent-> interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-c.Done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (c *Client) Close() {
	c.Close()
}

func NewClient(interrupt chan os.Signal) *Client {
	var u = url.URL{Scheme: "ws", Host: remoteAddress, Path: ""}
	if GetOutboundIP().Equal(localIP) {
		u = url.URL{Scheme: "ws", Host: localAddress, Path: ""}
	}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return nil
	}

	client := &Client{
		Conn:      c,
		Send:      make(chan []byte, 256),
		Done:      make(chan struct{}),
		Interrupt: interrupt,
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.onEvent()
	go client.onMessage()

	return client
}
