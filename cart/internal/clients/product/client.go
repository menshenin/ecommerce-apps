// Package product Структуры для работы с сервисом Product
package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"route256.ozon.ru/project/cart/internal/model"
	"route256.ozon.ru/project/cart/internal/pkg/errorgroup"
)

const getProductEndpoint = "/get_product"

type getProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type getProductRequest struct {
	Sku   uint32 `json:"sku"`
	Token string `json:"token"`
}

type httpRequestFunc func(r *http.Request) (*http.Response, error)

// Client Клиент для работы с сервисом товаров
type Client struct {
	url         string
	token       string
	limiter     *time.Ticker
	requestFunc httpRequestFunc
}

// Middleware Middleware
type Middleware func(httpRequestFunc) httpRequestFunc

// ClientOptionFunc Опции для конфигурации клиента
type ClientOptionFunc func(*Client)

// ErrProductServiceRequest Ошибка запроса к сервису Product
var ErrProductServiceRequest = errors.New("product service request error (non-200 status)")

// WithMiddleware Добавление middleware для запросов к сервису
func WithMiddleware(middleware Middleware) ClientOptionFunc {
	return func(client *Client) {
		client.requestFunc = middleware(client.requestFunc)
	}
}

// WithRPSLimit Ограничивает число одновременных запросов к сервису product
func WithRPSLimit(limit int64) ClientOptionFunc {
	return func(client *Client) {
		client.limiter = time.NewTicker(time.Second / time.Duration(limit))
	}
}

// WithRetryMiddleware Добавление middleware, реализующего ретрай запроса при 420/429 статусах
func WithRetryMiddleware(count int) ClientOptionFunc {
	return WithMiddleware(func(next httpRequestFunc) httpRequestFunc {
		return func(r *http.Request) (*http.Response, error) {
			if count < 2 {
				return next(r)
			}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			r.Body = io.NopCloser(bytes.NewReader(body))
			tries := 1
			for {
				resp, err := next(r)
				if tries < count && resp != nil && (resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == 420) {
					tries++
					r.Body = io.NopCloser(bytes.NewReader(body))
				} else {
					return resp, err
				}
			}
		}
	})
}

// WithToken Токен, передаваемый в сервис вместе с запросом
func WithToken(token string) ClientOptionFunc {
	return func(client *Client) {
		client.token = token
	}
}

// WithHTTPClient HTTP-client
func WithHTTPClient(httpClient *http.Client) ClientOptionFunc {
	return func(client *Client) {
		client.requestFunc = httpClient.Do
	}
}

// New Конструктор клиента для работы с сервисом товаров
func New(url string, opts ...ClientOptionFunc) *Client {
	client := &Client{url: url}

	for _, opt := range opts {
		opt(client)
	}
	if client.requestFunc == nil {
		httpClient := &http.Client{}
		client.requestFunc = httpClient.Do
	}

	return client
}

// GetItemsBySKU Получение информации о товарах
func (c *Client) GetItemsBySKU(ctx context.Context, skus ...model.SKU) (map[model.SKU]*model.Item, error) {
	results := make(map[model.SKU]*model.Item, len(skus))
	resultMutex := &sync.Mutex{}
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	g := errorgroup.New()
	for _, sku := range skus {
		g.Go(func(innerSku model.SKU) func() error {
			return func() error {
				if c.limiter != nil {
					select {
					case <-ctx.Done():
						return nil
					case <-c.limiter.C:
					}
				}
				item, err := c.requestItemBySKU(ctx, innerSku)
				if err != nil {
					return err
				}
				resultMutex.Lock()
				defer resultMutex.Unlock()
				results[innerSku] = item
				return nil
			}
		}(sku))
	}

	err := g.Wait()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (c *Client) requestItemBySKU(ctx context.Context, sku model.SKU) (*model.Item, error) {
	data := getProductRequest{
		Sku:   uint32(sku),
		Token: c.token,
	}
	byteData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url+getProductEndpoint, bytes.NewBuffer(byteData))
	if err != nil {
		return nil, err
	}
	resp, err := c.requestFunc(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrProductServiceRequest
	}
	decoder := json.NewDecoder(resp.Body)
	var result getProductResponse
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &model.Item{
		SKU:   sku,
		Name:  result.Name,
		Price: result.Price,
	}, nil
}
