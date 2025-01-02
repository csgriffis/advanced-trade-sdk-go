/**
 * Copyright 2024-present Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package test

import (
	"encoding/json"
	"os"

	"github.com/coinbase-samples/advanced-trade-sdk-go/client"
	"github.com/coinbase-samples/advanced-trade-sdk-go/credentials"
)

func setupClient() (client.RestClient, error) {
	credentials := &credentials.Credentials{}
	if err := json.Unmarshal([]byte(os.Getenv("ADV_CREDENTIALS")), credentials); err != nil {
		return nil, err
	}

	httpClient, err := client.DefaultHttpClient()
	if err != nil {
		return nil, err
	}
	restClient := client.NewRestClient(credentials, httpClient)
	return restClient, nil
}

func setupSocketClient() client.WebSocketClient {
	cfg := client.DefaultDialerConfig(client.DefaultMarketDataEndpoint)
	socketClient := client.NewWebSocketClient(cfg)
	return socketClient
}
