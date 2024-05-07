package product

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"route256.ozon.ru/project/cart/internal/model"
)

func TestClient_requestItemBySKU(t *testing.T) {
	type fields struct {
		url         string
		token       string
		requestFunc httpRequestFunc
	}
	type args struct {
		ctx context.Context
		sku model.SKU
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Item
		wantErr bool
	}{
		{
			name: "get item",
			fields: fields{
				url:   "test",
				token: "test",
				requestFunc: func(_ *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(`{"name": "testname", "price": 12}`)),
					}, nil
				},
			},
			args: args{
				ctx: context.Background(),
				sku: 123,
			},
			want: &model.Item{
				SKU:   123,
				Name:  "testname",
				Price: 12,
			},
			wantErr: false,
		},
		{
			name: "get item, non-200 code error",
			fields: fields{
				url:   "test",
				token: "test",
				requestFunc: func(_ *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(nil),
					}, nil
				},
			},
			args: args{
				ctx: context.Background(),
				sku: 123,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "get item, get error with empty response",
			fields: fields{
				url:   "test",
				token: "test",
				requestFunc: func(_ *http.Request) (*http.Response, error) {
					return nil, errors.New("some error")
				},
			},
			args: args{
				ctx: context.Background(),
				sku: 123,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				url:         tt.fields.url,
				token:       tt.fields.token,
				requestFunc: tt.fields.requestFunc,
			}
			got, err := c.requestItemBySKU(tt.args.ctx, tt.args.sku)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItemsBySKU() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requestItemBySKU() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithRetryMiddleware(t *testing.T) {
	client := New("test")
	callCount := 0
	client.requestFunc = func(_ *http.Request) (*http.Response, error) {
		callCount++
		return &http.Response{
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(nil),
		}, nil
	}
	WithRetryMiddleware(3)(client)
	_, err := client.requestItemBySKU(context.Background(), model.SKU(123))
	assert.ErrorIs(t, err, ErrProductServiceRequest)
	assert.Equal(t, 3, callCount)
}

func TestWithToken(t *testing.T) {
	client := New("test", WithToken("testtoken"))
	assert.Equal(t, client.token, "testtoken")
}

func TestClient_GetItemsBySKU(t *testing.T) {
	t.Run("get items, normal", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := New("test")
		client.requestFunc = func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"name": "testname", "price": 12}`)),
			}, nil
		}
		res, err := client.GetItemsBySKU(context.Background(), model.SKU(11), model.SKU(22))
		assert.NoError(t, err)
		assert.Equal(t, map[model.SKU]*model.Item{
			11: {
				SKU:   11,
				Name:  "testname",
				Price: 12,
			},
			22: {
				SKU:   22,
				Name:  "testname",
				Price: 12,
			},
		}, res)
	})

	t.Run("get items, error", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		client := New("test")
		client.requestFunc = func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(nil),
			}, nil
		}
		res, err := client.GetItemsBySKU(context.Background(), model.SKU(11), model.SKU(22))
		assert.Error(t, err)
		assert.Empty(t, res)
	})
}

func TestWithRPSLimit(t *testing.T) {
	defer goleak.VerifyNone(t)
	client := New("test", WithRPSLimit(2))

	sem := make(chan struct{})
	defer func() {
		close(sem)
	}()
	requestCount := atomic.Int32{}
	wg := sync.WaitGroup{}
	client.requestFunc = func(_ *http.Request) (*http.Response, error) {
		requestCount.Add(1)
		wg.Done()
		<-sem
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"name": "testname", "price": 12}`)),
		}, nil
	}
	wg.Add(2)
	go func() {
		_, _ = client.GetItemsBySKU(context.Background(), model.SKU(11), model.SKU(22), model.SKU(33))
	}()
	wg.Wait()
	assert.EqualValues(t, 2, requestCount.Load()) // выполняются 2 запроса, третий ждёт
	wg.Add(1)
	sem <- struct{}{}
	sem <- struct{}{}
	wg.Wait()
	assert.EqualValues(t, 3, requestCount.Load()) // выполняется третий запрос
	sem <- struct{}{}
	// все запросы выполнены, горутины завершены
}
