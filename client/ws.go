/*
Copyright © 2025 Chris Griffis <dev@chrisgriffis.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coinbase-samples/advanced-trade-sdk-go/model"
	"github.com/coinbase-samples/core-go"

	"github.com/gorilla/websocket"
)

var DefaultMarketDataEndpoint = "wss://advanced-trade-ws.coinbase.com"
var DefaultUserOrderDataEndpoint = "wss://advanced-trade-ws-user.coinbase.com"

type WebSocketClient interface {
	// Dial establishes the WebSocket connection using the underlying
	// DialerConfig.
	Dial(ctx context.Context) error

	// Close closes the WebSocket connection.
	Close() error

	// Subscribe sends a subscription request for one or more
	// channels (e.g., "ticker", "level2", "user", etc.).
	Subscribe(request model.SocketRequest) error

	// Listen continuously reads messages from the WebSocket
	// connection, passing each message to the provided handler.
	// It blocks until the context is cancelled or an error occurs.
	Listen(ctx context.Context, handleMessage core.OnWebSocketTextMessage) <-chan error
}

func NewWebSocketClient(cfg core.DialerConfig) WebSocketClient {
	return &webSocketImpl{config: cfg}
}

type webSocketImpl struct {
	config core.DialerConfig
	conn   *core.WebSocketConnection
}

// Dial uses core.DialWebSocket to establish the underlying WebSocket connection.
func (w *webSocketImpl) Dial(ctx context.Context) error {
	var err error
	w.conn, err = core.DialWebSocket(ctx, w.config)
	if err != nil {
		return fmt.Errorf("error dialing web socket: %v", err)
	}

	return nil
}

// Close gracefully closes the WebSocket if it is open.
func (w *webSocketImpl) Close() error {
	if w.conn == nil {
		return nil
	}
	return w.conn.Close() // conn.Conn is the *websocket.Conn from gorilla/websocket
}

// Subscribe sends a "subscribe" message to the WebSocket for the given channels.
func (w *webSocketImpl) Subscribe(request model.SocketRequest) error {
	if w.conn == nil {
		return fmt.Errorf("cannot subscribe: no active WebSocket connection")
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal subscribe message: %v", err)
	}

	// Use the wrapped *websocket.Conn to write the message.
	if err = w.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to write subscribe message: %v", err)
	}
	return nil
}

// Listen accepts a message handler and returns an error channel. If an error occurs setting up the listener, the
// error is returned on the channel.
func (w *webSocketImpl) Listen(ctx context.Context, handleMessage core.OnWebSocketTextMessage) <-chan error {
	var errChan = make(chan error)

	if w.conn == nil {
		errChan <- fmt.Errorf("cannot listen: no active WebSocket connection")
		return errChan
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				errChan <- fmt.Errorf("context done: %w", ctx.Err())
				return
			default:
				// Read next message from the gorilla/websocket connection.
				msgType, message, err := w.conn.ReadMessage()
				if err != nil {
					errChan <- fmt.Errorf("failed to read message: %w", err)
					// do not return from the goroutine, we want the connection to stay open
					continue
				}

				// You can check message types if needed (e.g., TextMessage vs. BinaryMessage).
				if msgType != websocket.TextMessage {
					// Possibly handle or ignore non-text messages here.
					continue
				}
				handleMessage(message)
			}
		}
	}()

	return errChan
}

func DefaultDialerConfig(url string) core.DialerConfig { return core.DefaultDialerConfig(url) }
