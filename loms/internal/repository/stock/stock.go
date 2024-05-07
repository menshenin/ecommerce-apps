package stock

import (
	"sync"

	"route256.ozon.ru/project/loms/internal/model"
)

// Stock Реализация склада, хранящего все товары в памяти
type Stock struct {
	mutex   *sync.RWMutex
	storage map[model.SKU]model.StockItem
}

// WriteOff Списывание зарезервированных товаров
func (s Stock) WriteOff(sku model.SKU, count uint16) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if si, ok := s.storage[sku]; ok {
		si.Reserved -= int32(count)
		si.TotalCount -= int32(count)
		s.storage[sku] = si
		return nil
	}

	return model.ErrSkuNotFound
}

// Reserve Резервирование товаров
func (s Stock) Reserve(sku model.SKU, count uint16) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if si, ok := s.storage[sku]; ok {
		reserved := si.Reserved + int32(count)
		if reserved > si.TotalCount {
			return model.ErrReserveMoreThenTotalCount
		}
		si.Reserved = reserved
		s.storage[sku] = si
		return nil
	}

	return model.ErrSkuNotFound
}

// AvailableCount Доступное количество товара
func (s Stock) AvailableCount(sku model.SKU) (int32, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if si, ok := s.storage[sku]; ok {
		return si.TotalCount - si.Reserved, nil
	}

	return 0, model.ErrSkuNotFound
}

// CancelReserve Отмена резервирования
func (s Stock) CancelReserve(sku model.SKU, count uint16) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if si, ok := s.storage[sku]; ok {
		if count > uint16(si.Reserved) {
			return model.ErrCancelMoreThenReserved
		}
		si.Reserved -= int32(count)
		s.storage[sku] = si
		return nil
	}

	return model.ErrSkuNotFound
}

// Load Загрузка данных по товарам
func (s Stock) Load(items []model.StockItem) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, item := range items {
		s.storage[item.SKU] = item
	}
}
