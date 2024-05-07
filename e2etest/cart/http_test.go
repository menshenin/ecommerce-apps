package cart

import (
	_ "embed"
	"io"
	"net/http"
	"strings"
	"testing"
)

const baseUrl = "http://localhost:8082"

type httpTest struct {
	name         string
	method       string
	path         string
	body         io.Reader
	expectedCode int
	expectedBody string
}

func TestE2E(t *testing.T) {
	tests := []httpTest{
		{
			name:         "test clear",
			method:       http.MethodDelete,
			path:         "/user/31337/cart",
			expectedCode: 204,
		},
		{
			name:         "check empty cart",
			method:       http.MethodGet,
			path:         "/user/31337/cart",
			expectedCode: 404,
		},
		{
			name:         "add item to cart, normal",
			method:       http.MethodPost,
			path:         "/user/31337/cart/773297411",
			body:         strings.NewReader(`{"count": 10}`),
			expectedCode: 200,
		},
		{
			name:         "check cart again, expect 773297411 sku",
			method:       http.MethodGet,
			path:         "/user/31337/cart",
			expectedCode: 200,
			expectedBody: `{"items":[{"sku_id":773297411,"name":"Кроссовки Nike JORDAN","count":10,"price":2202}],"total_price":22020}`,
		},
		{
			name:         "remove sku from cart",
			method:       http.MethodDelete,
			path:         "/user/31337/cart/773297411",
			expectedCode: 204,
		},
		{
			name:         "add item to cart",
			method:       http.MethodPost,
			path:         "/user/31337/cart/773297411",
			body:         strings.NewReader(`{"count": 5}`),
			expectedCode: 200,
		},
		{
			name:         "checkout",
			method:       http.MethodPost,
			path:         "/cart/checkout",
			expectedCode: 200,
			body:         strings.NewReader(`{"user_id": 31337}`),
			expectedBody: `{"order_id":1}`,
		},
		{
			name:         "check cart is empty",
			method:       http.MethodGet,
			path:         "/user/31337/cart",
			expectedCode: 404,
		},
	}
	client := http.Client{}
	for _, test := range tests {
		t.Run(test.name, runTestFunc(&client, test))
	}
}

func runTestFunc(client *http.Client, test httpTest) func(t *testing.T) {
	return func(t *testing.T) {
		req, err := http.NewRequest(test.method, baseUrl+test.path, test.body)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != test.expectedCode {
			t.Fatalf("want code: %d, get: %d", test.expectedCode, resp.StatusCode)
		}
		if test.expectedBody != "" {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read body: %v", err)
			}
			strBody := strings.TrimSpace(string(body))
			if strBody != test.expectedBody {
				t.Fatalf("invalid body = got %v, want %v", strBody, test.expectedBody)
			}
		}
	}
}
