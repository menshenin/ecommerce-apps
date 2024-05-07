package inmemorycart

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"route256.ozon.ru/project/cart/internal/model"
)

func TestRepository_GetByUserID(t *testing.T) {
	t.Parallel()
	r := New()
	t.Run("get not existed cart", func(t *testing.T) {
		cart, err := r.GetByUserID(context.Background(), model.UserID(123))
		assert.Nil(t, cart)
		assert.ErrorIs(t, err, model.ErrCartNotFound)
	})
	t.Run("create and get cart", func(t *testing.T) {
		_, _ = r.Create(context.Background(), model.UserID(123))
		cart, err := r.GetByUserID(context.Background(), model.UserID(123))
		assert.NoError(t, err)
		assert.Equal(t, cart.UserID, model.UserID(123))
	})
}

func BenchmarkRepository_Create(b *testing.B) {
	r := New()
	for i := 0; i < b.N; i++ {
		_, err := r.Create(context.Background(), model.UserID(i))
		if err != nil {
			b.Fatal(err)
		}
	}
}
