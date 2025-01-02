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

package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/coinbase-samples/advanced-trade-sdk-go/model"
)

func TestListen(t *testing.T) {
	client := setupSocketClient()
	defer client.Close()

	err := client.Dial(context.Background())
	if err != nil {
		t.Errorf("Failed to dial: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageC := make(chan []byte)
	errC := client.Listen(ctx, func(message []byte) {
		messageC <- message
	})

	if err := client.Subscribe(model.SocketRequest{Type: "subscribe", Channel: "market_trades", ProductIds: []string{"BTC-USD"}}); err != nil {
		t.Errorf("Failed to subscribe: %v", err)
	}

	for {
		select {
		case message := <-messageC:
			fmt.Println(string(message))
		case err := <-errC:
			t.Errorf(err.Error())
		case <-ctx.Done():
			return
		default:
		}
	}
}
