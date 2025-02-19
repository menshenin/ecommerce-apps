// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package mocks

//go:generate minimock -i route256.ozon.ru/project/cart/internal/service/cart.ItemRepository -o item_repository_mock.go -n ItemRepositoryMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/cart/internal/model"
)

// ItemRepositoryMock implements cart.ItemRepository
type ItemRepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetItemsBySKU          func(ctx context.Context, sku ...model.SKU) (m1 map[model.SKU]*model.Item, err error)
	inspectFuncGetItemsBySKU   func(ctx context.Context, sku ...model.SKU)
	afterGetItemsBySKUCounter  uint64
	beforeGetItemsBySKUCounter uint64
	GetItemsBySKUMock          mItemRepositoryMockGetItemsBySKU
}

// NewItemRepositoryMock returns a mock for cart.ItemRepository
func NewItemRepositoryMock(t minimock.Tester) *ItemRepositoryMock {
	m := &ItemRepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetItemsBySKUMock = mItemRepositoryMockGetItemsBySKU{mock: m}
	m.GetItemsBySKUMock.callArgs = []*ItemRepositoryMockGetItemsBySKUParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mItemRepositoryMockGetItemsBySKU struct {
	mock               *ItemRepositoryMock
	defaultExpectation *ItemRepositoryMockGetItemsBySKUExpectation
	expectations       []*ItemRepositoryMockGetItemsBySKUExpectation

	callArgs []*ItemRepositoryMockGetItemsBySKUParams
	mutex    sync.RWMutex
}

// ItemRepositoryMockGetItemsBySKUExpectation specifies expectation struct of the ItemRepository.GetItemsBySKU
type ItemRepositoryMockGetItemsBySKUExpectation struct {
	mock    *ItemRepositoryMock
	params  *ItemRepositoryMockGetItemsBySKUParams
	results *ItemRepositoryMockGetItemsBySKUResults
	Counter uint64
}

// ItemRepositoryMockGetItemsBySKUParams contains parameters of the ItemRepository.GetItemsBySKU
type ItemRepositoryMockGetItemsBySKUParams struct {
	ctx context.Context
	sku []model.SKU
}

// ItemRepositoryMockGetItemsBySKUResults contains results of the ItemRepository.GetItemsBySKU
type ItemRepositoryMockGetItemsBySKUResults struct {
	m1  map[model.SKU]*model.Item
	err error
}

// Expect sets up expected params for ItemRepository.GetItemsBySKU
func (mmGetItemsBySKU *mItemRepositoryMockGetItemsBySKU) Expect(ctx context.Context, sku ...model.SKU) *mItemRepositoryMockGetItemsBySKU {
	if mmGetItemsBySKU.mock.funcGetItemsBySKU != nil {
		mmGetItemsBySKU.mock.t.Fatalf("ItemRepositoryMock.GetItemsBySKU mock is already set by Set")
	}

	if mmGetItemsBySKU.defaultExpectation == nil {
		mmGetItemsBySKU.defaultExpectation = &ItemRepositoryMockGetItemsBySKUExpectation{}
	}

	mmGetItemsBySKU.defaultExpectation.params = &ItemRepositoryMockGetItemsBySKUParams{ctx, sku}
	for _, e := range mmGetItemsBySKU.expectations {
		if minimock.Equal(e.params, mmGetItemsBySKU.defaultExpectation.params) {
			mmGetItemsBySKU.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetItemsBySKU.defaultExpectation.params)
		}
	}

	return mmGetItemsBySKU
}

// Inspect accepts an inspector function that has same arguments as the ItemRepository.GetItemsBySKU
func (mmGetItemsBySKU *mItemRepositoryMockGetItemsBySKU) Inspect(f func(ctx context.Context, sku ...model.SKU)) *mItemRepositoryMockGetItemsBySKU {
	if mmGetItemsBySKU.mock.inspectFuncGetItemsBySKU != nil {
		mmGetItemsBySKU.mock.t.Fatalf("Inspect function is already set for ItemRepositoryMock.GetItemsBySKU")
	}

	mmGetItemsBySKU.mock.inspectFuncGetItemsBySKU = f

	return mmGetItemsBySKU
}

// Return sets up results that will be returned by ItemRepository.GetItemsBySKU
func (mmGetItemsBySKU *mItemRepositoryMockGetItemsBySKU) Return(m1 map[model.SKU]*model.Item, err error) *ItemRepositoryMock {
	if mmGetItemsBySKU.mock.funcGetItemsBySKU != nil {
		mmGetItemsBySKU.mock.t.Fatalf("ItemRepositoryMock.GetItemsBySKU mock is already set by Set")
	}

	if mmGetItemsBySKU.defaultExpectation == nil {
		mmGetItemsBySKU.defaultExpectation = &ItemRepositoryMockGetItemsBySKUExpectation{mock: mmGetItemsBySKU.mock}
	}
	mmGetItemsBySKU.defaultExpectation.results = &ItemRepositoryMockGetItemsBySKUResults{m1, err}
	return mmGetItemsBySKU.mock
}

// Set uses given function f to mock the ItemRepository.GetItemsBySKU method
func (mmGetItemsBySKU *mItemRepositoryMockGetItemsBySKU) Set(f func(ctx context.Context, sku ...model.SKU) (m1 map[model.SKU]*model.Item, err error)) *ItemRepositoryMock {
	if mmGetItemsBySKU.defaultExpectation != nil {
		mmGetItemsBySKU.mock.t.Fatalf("Default expectation is already set for the ItemRepository.GetItemsBySKU method")
	}

	if len(mmGetItemsBySKU.expectations) > 0 {
		mmGetItemsBySKU.mock.t.Fatalf("Some expectations are already set for the ItemRepository.GetItemsBySKU method")
	}

	mmGetItemsBySKU.mock.funcGetItemsBySKU = f
	return mmGetItemsBySKU.mock
}

// When sets expectation for the ItemRepository.GetItemsBySKU which will trigger the result defined by the following
// Then helper
func (mmGetItemsBySKU *mItemRepositoryMockGetItemsBySKU) When(ctx context.Context, sku ...model.SKU) *ItemRepositoryMockGetItemsBySKUExpectation {
	if mmGetItemsBySKU.mock.funcGetItemsBySKU != nil {
		mmGetItemsBySKU.mock.t.Fatalf("ItemRepositoryMock.GetItemsBySKU mock is already set by Set")
	}

	expectation := &ItemRepositoryMockGetItemsBySKUExpectation{
		mock:   mmGetItemsBySKU.mock,
		params: &ItemRepositoryMockGetItemsBySKUParams{ctx, sku},
	}
	mmGetItemsBySKU.expectations = append(mmGetItemsBySKU.expectations, expectation)
	return expectation
}

// Then sets up ItemRepository.GetItemsBySKU return parameters for the expectation previously defined by the When method
func (e *ItemRepositoryMockGetItemsBySKUExpectation) Then(m1 map[model.SKU]*model.Item, err error) *ItemRepositoryMock {
	e.results = &ItemRepositoryMockGetItemsBySKUResults{m1, err}
	return e.mock
}

// GetItemsBySKU implements cart.ItemRepository
func (mmGetItemsBySKU *ItemRepositoryMock) GetItemsBySKU(ctx context.Context, sku ...model.SKU) (m1 map[model.SKU]*model.Item, err error) {
	mm_atomic.AddUint64(&mmGetItemsBySKU.beforeGetItemsBySKUCounter, 1)
	defer mm_atomic.AddUint64(&mmGetItemsBySKU.afterGetItemsBySKUCounter, 1)

	if mmGetItemsBySKU.inspectFuncGetItemsBySKU != nil {
		mmGetItemsBySKU.inspectFuncGetItemsBySKU(ctx, sku...)
	}

	mm_params := ItemRepositoryMockGetItemsBySKUParams{ctx, sku}

	// Record call args
	mmGetItemsBySKU.GetItemsBySKUMock.mutex.Lock()
	mmGetItemsBySKU.GetItemsBySKUMock.callArgs = append(mmGetItemsBySKU.GetItemsBySKUMock.callArgs, &mm_params)
	mmGetItemsBySKU.GetItemsBySKUMock.mutex.Unlock()

	for _, e := range mmGetItemsBySKU.GetItemsBySKUMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.m1, e.results.err
		}
	}

	if mmGetItemsBySKU.GetItemsBySKUMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetItemsBySKU.GetItemsBySKUMock.defaultExpectation.Counter, 1)
		mm_want := mmGetItemsBySKU.GetItemsBySKUMock.defaultExpectation.params
		mm_got := ItemRepositoryMockGetItemsBySKUParams{ctx, sku}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetItemsBySKU.t.Errorf("ItemRepositoryMock.GetItemsBySKU got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetItemsBySKU.GetItemsBySKUMock.defaultExpectation.results
		if mm_results == nil {
			mmGetItemsBySKU.t.Fatal("No results are set for the ItemRepositoryMock.GetItemsBySKU")
		}
		return (*mm_results).m1, (*mm_results).err
	}
	if mmGetItemsBySKU.funcGetItemsBySKU != nil {
		return mmGetItemsBySKU.funcGetItemsBySKU(ctx, sku...)
	}
	mmGetItemsBySKU.t.Fatalf("Unexpected call to ItemRepositoryMock.GetItemsBySKU. %v %v", ctx, sku)
	return
}

// GetItemsBySKUAfterCounter returns a count of finished ItemRepositoryMock.GetItemsBySKU invocations
func (mmGetItemsBySKU *ItemRepositoryMock) GetItemsBySKUAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetItemsBySKU.afterGetItemsBySKUCounter)
}

// GetItemsBySKUBeforeCounter returns a count of ItemRepositoryMock.GetItemsBySKU invocations
func (mmGetItemsBySKU *ItemRepositoryMock) GetItemsBySKUBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetItemsBySKU.beforeGetItemsBySKUCounter)
}

// Calls returns a list of arguments used in each call to ItemRepositoryMock.GetItemsBySKU.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetItemsBySKU *mItemRepositoryMockGetItemsBySKU) Calls() []*ItemRepositoryMockGetItemsBySKUParams {
	mmGetItemsBySKU.mutex.RLock()

	argCopy := make([]*ItemRepositoryMockGetItemsBySKUParams, len(mmGetItemsBySKU.callArgs))
	copy(argCopy, mmGetItemsBySKU.callArgs)

	mmGetItemsBySKU.mutex.RUnlock()

	return argCopy
}

// MinimockGetItemsBySKUDone returns true if the count of the GetItemsBySKU invocations corresponds
// the number of defined expectations
func (m *ItemRepositoryMock) MinimockGetItemsBySKUDone() bool {
	for _, e := range m.GetItemsBySKUMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetItemsBySKUMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetItemsBySKUCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetItemsBySKU != nil && mm_atomic.LoadUint64(&m.afterGetItemsBySKUCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetItemsBySKUInspect logs each unmet expectation
func (m *ItemRepositoryMock) MinimockGetItemsBySKUInspect() {
	for _, e := range m.GetItemsBySKUMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ItemRepositoryMock.GetItemsBySKU with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetItemsBySKUMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetItemsBySKUCounter) < 1 {
		if m.GetItemsBySKUMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ItemRepositoryMock.GetItemsBySKU")
		} else {
			m.t.Errorf("Expected call to ItemRepositoryMock.GetItemsBySKU with params: %#v", *m.GetItemsBySKUMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetItemsBySKU != nil && mm_atomic.LoadUint64(&m.afterGetItemsBySKUCounter) < 1 {
		m.t.Error("Expected call to ItemRepositoryMock.GetItemsBySKU")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ItemRepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetItemsBySKUInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ItemRepositoryMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ItemRepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetItemsBySKUDone()
}
