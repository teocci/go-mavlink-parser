// Package wsnet
// Created by RTT.
// Author: teocci@yandex.com on 2021-Sep-30
package wsnet

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/teocci/go-mavlink-parser/src/datamgr"
)

const (
	defaultServerIP = "localhost"
	wsPort   = 7474

	CMDPing     = "ping"
	CMDPong     = "pong"
	CMDRegister = "register"

	CMDWebsocketConnected = "websocket-connected"
	CMDConnectServices    = "connect-services"
	CMDUpdateTelemetry = "update-telemetry"

	RoleWebConsumer     = "web-consumer"
	RoleTelemetryPusher = "telemetry-pusher"
	RoleStreamingPusher = "streaming-pusher"

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

	conf datamgr.InitConf
)

// Client is a middleman between the websocket server and this application
type Client struct {
	// The websocket connection.
	Conn *websocket.Conn

	ConnectionID int64
	WorkerID int64

	// Buffered channel of outbound messages.
	Send      chan []byte
	Done      chan struct{}
	Interrupt chan os.Signal
}

// onMessage reads messages from the websocket connection
//
// Runs onMessage in a per-connection goroutine. The application ensures
// that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) onMessage() {
	defer func() {
		_ = c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Fatalf("error: %v", err)
			return
		}

		var dat map[string]interface{}
		if err := json.Unmarshal(message, &dat); err != nil {
			panic(err)
		}

		switch dat["cmd"] {
		case CMDWebsocketConnected:
			// {"cmd":"ws-connected","connection_id":xxx}
			c.ConnectionID = int64(dat["connection_id"].(float64))
			c.WorkerID = int64(dat["worker_id"].(float64))

			req := &datamgr.Register{
				CMD:       CMDRegister,
				ConnID:    c.ConnectionID,
				WorkerID:    c.WorkerID,
				ModuleTag: conf.ModuleTag,
				DroneID:   conf.DroneID,
				Role:      RoleTelemetryPusher,
			}

			jsonData, err := json.Marshal(req)
			if err != nil {
				log.Fatalf("onMessage(): %v", err)
			}
			c.Send <- jsonData
		case CMDUpdateTelemetry:
			continue
		default:
			log.Printf("recv: %s", message)
		}
	}
}

// onEvent writes messages from the hub to the websocket connection.
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
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)
			//data := string(message[:])
			//log.Println("onEvent-> message")
			//log.Printf("%#v", message)
			//log.Printf("%#v", data)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
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

func NewClient(c datamgr.InitConf) *Client {
	conf = c
	WebsocketServerIP(conf.WSHost)
	var u = url.URL{Scheme: "ws", Host: serverAddress, Path: ""}
	log.Printf("connecting to %s", u.String())

	wsConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return nil
	}

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	client := &Client{
		Conn:      wsConn,
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
