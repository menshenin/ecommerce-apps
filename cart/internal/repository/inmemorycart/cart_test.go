package inmemorycart

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/cart/internal/model"
)

func TestCart_AddItemBySKU(t *testing.T) {
	t.Run("add item, simple", func(t *testing.T) {
		type fields struct {
			items map[model.SKU]model.CartItem
		}
		type args struct {
			sku   model.SKU
			count uint16
		}
		tests := []struct {
			name      string
			fields    fields
			args      args
			wantItems map[model.SKU]model.CartItem
		}{
			{
				name: "add new item",
				fields: fields{
					items: map[model.SKU]model.CartItem{
						123: {
							SKU:   123,
							Count: 1,
						},
					},
				},
				args: args{
					sku:   456,
					count: 3,
				},
				wantItems: map[model.SKU]model.CartItem{
					123: {
						SKU:   123,
						Count: 1,
					},
					456: {
						SKU:   456,
						Count: 3,
					},
				},
			},
			{
				name: "add existed item",
				fields: fields{
					items: map[model.SKU]model.CartItem{
						123: {
							SKU:   123,
							Count: 1,
						},
					},
				},
				args: args{
					sku:   123,
					count: 3,
				},
				wantItems: map[model.SKU]model.CartItem{
					123: {
						SKU:   123,
						Count: 4,
					},
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				c := &Cart{
					RWMutex: &sync.RWMutex{},
					items:   tt.fields.items,
				}
				assert.NoError(t, c.AddItemBySKU(tt.args.sku, tt.args.count))
				if !reflect.DeepEqual(c.items, tt.wantItems) {
					t.Errorf("Invalid items got = %v, want %v", c.items, tt.wantItems)
				}
			})
		}
	})

	t.Run("add item, concurrent", func(t *testing.T) {
		cart := NewCart(model.UserID(1))
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				err := cart.AddItemBySKU(model.SKU(idx%10), 10)
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()
		assert.Len(t, cart.items, 10)
		for _, item := range cart.items {
			assert.EqualValues(t, 100, item.Count)
		}
	})
}

func TestCart_Clear(t *testing.T) {
	type fields struct {
		items map[model.SKU]model.CartItem
	}
	tests := []struct {
		name      string
		fields    fields
		wantItems map[model.SKU]model.CartItem
	}{
		{
			name: "clear existed items",
			fields: fields{
				items: map[model.SKU]model.CartItem{
					123: {
						SKU:   123,
						Count: 1,
					},
				},
			},
			wantItems: map[model.SKU]model.CartItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cart{
				RWMutex: &sync.RWMutex{},
				items:   tt.fields.items,
			}
			assert.NoError(t, c.Clear())
			if !reflect.DeepEqual(c.items, tt.wantItems) {
				t.Errorf("Invalid items got = %v, want %v", c.items, tt.wantItems)
			}
		})
	}
}

func TestCart_DeleteItemBySKU(t *testing.T) {
	type fields struct {
		items map[model.SKU]model.CartItem
	}
	type args struct {
		sku model.SKU
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantItems map[model.SKU]model.CartItem
	}{
		{
			name: "delete existed item",
			fields: fields{
				items: map[model.SKU]model.CartItem{
					123: {
						SKU:   123,
						Count: 5,
					},
				},
			},
			args:      args{sku: 123},
			wantItems: map[model.SKU]model.CartItem{},
		},
		{
			name: "delete not existed item",
			fields: fields{
				items: map[model.SKU]model.CartItem{
					123: {
						SKU:   123,
						Count: 5,
					},
				},
			},
			args: args{sku: 456},
			wantItems: map[model.SKU]model.CartItem{
				123: {
					SKU:   123,
					Count: 5,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cart{
				RWMutex: &sync.RWMutex{},
				items:   tt.fields.items,
			}
			assert.NoError(t, c.DeleteItemBySKU(tt.args.sku))
			if !reflect.DeepEqual(c.items, tt.wantItems) {
				t.Errorf("Invalid items got = %v, want %v", c.items, tt.wantItems)
			}
		})
	}
}

func TestCart_GetItems(t *testing.T) {
	type fields struct {
		items map[model.SKU]model.CartItem
	}
	tests := []struct {
		name   string
		fields fields
		want   []model.CartItem
	}{
		{
			name: "get existed items sorted by SKU",
			fields: fields{
				items: map[model.SKU]model.CartItem{
					456: {
						SKU:   456,
						Count: 5,
					},
					123: {
						SKU:   123,
						Count: 1,
					},
				},
			},
			want: []model.CartItem{
				{
					SKU:   456,
					Count: 5,
				},
				{
					SKU:   123,
					Count: 1,
				},
			},
		},
		{
			name:   "get empty",
			fields: fields{},
			want:   []model.CartItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cart{
				RWMutex: &sync.RWMutex{},
				items:   tt.fields.items,
			}
			got, err := c.GetItems()
			assert.NoError(t, err)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() got = %v, want %v", got, tt.want)
			}
		})
	}
}
