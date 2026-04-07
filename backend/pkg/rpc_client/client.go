package rpc_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ybbus/jsonrpc/v3"
)

func CallJsonRpc[T any](rpcClient jsonrpc.RPCClient, ctx context.Context, method string, args ...any) (T, error) {
	var result T
	res, err := rpcClient.Call(ctx, method, args...)
	if err != nil {
		return result, fmt.Errorf("call rpc %s: %w", method, err)
	}
	err = res.GetObject(&result)
	if err != nil {
		return result, fmt.Errorf("get rpc result: %w", err)
	}
	return result, nil
}

type RestRpcClient struct {
	client  http.Client
	BaseUrl string
}

func NewRestRpcClient(baseUrl string) *RestRpcClient {
	return &RestRpcClient{
		client:  http.Client{},
		BaseUrl: baseUrl,
	}
}

func (c *RestRpcClient) Call(
	ctx context.Context,
	namespace string,
	method string,
	result any,
	params ...any,
) error {
	payload := map[string]any{
		"method": fmt.Sprintf("%s_%s", namespace, method),
		"params": params,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseUrl,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(result)

}
