// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package mocks

//go:generate minimock -i route256.ozon.ru/project/loms/internal/service/order.OrdersRepository -o orders_repository_mock.go -n OrdersRepositoryMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/loms/internal/model"
)

// OrdersRepositoryMock implements order.OrdersRepository
type OrdersRepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcCreate          func(ctx context.Context, userID model.UserID, items []model.OrderItem) (op1 *model.Order, err error)
	inspectFuncCreate   func(ctx context.Context, userID model.UserID, items []model.OrderItem)
	afterCreateCounter  uint64
	beforeCreateCounter uint64
	CreateMock          mOrdersRepositoryMockCreate

	funcGetByID          func(ctx context.Context, id model.OrderID) (op1 *model.Order, err error)
	inspectFuncGetByID   func(ctx context.Context, id model.OrderID)
	afterGetByIDCounter  uint64
	beforeGetByIDCounter uint64
	GetByIDMock          mOrdersRepositoryMockGetByID
}

// NewOrdersRepositoryMock returns a mock for order.OrdersRepository
func NewOrdersRepositoryMock(t minimock.Tester) *OrdersRepositoryMock {
	m := &OrdersRepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateMock = mOrdersRepositoryMockCreate{mock: m}
	m.CreateMock.callArgs = []*OrdersRepositoryMockCreateParams{}

	m.GetByIDMock = mOrdersRepositoryMockGetByID{mock: m}
	m.GetByIDMock.callArgs = []*OrdersRepositoryMockGetByIDParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mOrdersRepositoryMockCreate struct {
	mock               *OrdersRepositoryMock
	defaultExpectation *OrdersRepositoryMockCreateExpectation
	expectations       []*OrdersRepositoryMockCreateExpectation

	callArgs []*OrdersRepositoryMockCreateParams
	mutex    sync.RWMutex
}

// OrdersRepositoryMockCreateExpectation specifies expectation struct of the OrdersRepository.Create
type OrdersRepositoryMockCreateExpectation struct {
	mock    *OrdersRepositoryMock
	params  *OrdersRepositoryMockCreateParams
	results *OrdersRepositoryMockCreateResults
	Counter uint64
}

// OrdersRepositoryMockCreateParams contains parameters of the OrdersRepository.Create
type OrdersRepositoryMockCreateParams struct {
	ctx    context.Context
	userID model.UserID
	items  []model.OrderItem
}

// OrdersRepositoryMockCreateResults contains results of the OrdersRepository.Create
type OrdersRepositoryMockCreateResults struct {
	op1 *model.Order
	err error
}

// Expect sets up expected params for OrdersRepository.Create
func (mmCreate *mOrdersRepositoryMockCreate) Expect(ctx context.Context, userID model.UserID, items []model.OrderItem) *mOrdersRepositoryMockCreate {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("OrdersRepositoryMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &OrdersRepositoryMockCreateExpectation{}
	}

	mmCreate.defaultExpectation.params = &OrdersRepositoryMockCreateParams{ctx, userID, items}
	for _, e := range mmCreate.expectations {
		if minimock.Equal(e.params, mmCreate.defaultExpectation.params) {
			mmCreate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmCreate.defaultExpectation.params)
		}
	}

	return mmCreate
}

// Inspect accepts an inspector function that has same arguments as the OrdersRepository.Create
func (mmCreate *mOrdersRepositoryMockCreate) Inspect(f func(ctx context.Context, userID model.UserID, items []model.OrderItem)) *mOrdersRepositoryMockCreate {
	if mmCreate.mock.inspectFuncCreate != nil {
		mmCreate.mock.t.Fatalf("Inspect function is already set for OrdersRepositoryMock.Create")
	}

	mmCreate.mock.inspectFuncCreate = f

	return mmCreate
}

// Return sets up results that will be returned by OrdersRepository.Create
func (mmCreate *mOrdersRepositoryMockCreate) Return(op1 *model.Order, err error) *OrdersRepositoryMock {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("OrdersRepositoryMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &OrdersRepositoryMockCreateExpectation{mock: mmCreate.mock}
	}
	mmCreate.defaultExpectation.results = &OrdersRepositoryMockCreateResults{op1, err}
	return mmCreate.mock
}

// Set uses given function f to mock the OrdersRepository.Create method
func (mmCreate *mOrdersRepositoryMockCreate) Set(f func(ctx context.Context, userID model.UserID, items []model.OrderItem) (op1 *model.Order, err error)) *OrdersRepositoryMock {
	if mmCreate.defaultExpectation != nil {
		mmCreate.mock.t.Fatalf("Default expectation is already set for the OrdersRepository.Create method")
	}

	if len(mmCreate.expectations) > 0 {
		mmCreate.mock.t.Fatalf("Some expectations are already set for the OrdersRepository.Create method")
	}

	mmCreate.mock.funcCreate = f
	return mmCreate.mock
}

// When sets expectation for the OrdersRepository.Create which will trigger the result defined by the following
// Then helper
func (mmCreate *mOrdersRepositoryMockCreate) When(ctx context.Context, userID model.UserID, items []model.OrderItem) *OrdersRepositoryMockCreateExpectation {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("OrdersRepositoryMock.Create mock is already set by Set")
	}

	expectation := &OrdersRepositoryMockCreateExpectation{
		mock:   mmCreate.mock,
		params: &OrdersRepositoryMockCreateParams{ctx, userID, items},
	}
	mmCreate.expectations = append(mmCreate.expectations, expectation)
	return expectation
}

// Then sets up OrdersRepository.Create return parameters for the expectation previously defined by the When method
func (e *OrdersRepositoryMockCreateExpectation) Then(op1 *model.Order, err error) *OrdersRepositoryMock {
	e.results = &OrdersRepositoryMockCreateResults{op1, err}
	return e.mock
}

// Create implements order.OrdersRepository
func (mmCreate *OrdersRepositoryMock) Create(ctx context.Context, userID model.UserID, items []model.OrderItem) (op1 *model.Order, err error) {
	mm_atomic.AddUint64(&mmCreate.beforeCreateCounter, 1)
	defer mm_atomic.AddUint64(&mmCreate.afterCreateCounter, 1)

	if mmCreate.inspectFuncCreate != nil {
		mmCreate.inspectFuncCreate(ctx, userID, items)
	}

	mm_params := OrdersRepositoryMockCreateParams{ctx, userID, items}

	// Record call args
	mmCreate.CreateMock.mutex.Lock()
	mmCreate.CreateMock.callArgs = append(mmCreate.CreateMock.callArgs, &mm_params)
	mmCreate.CreateMock.mutex.Unlock()

	for _, e := range mmCreate.CreateMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.op1, e.results.err
		}
	}

	if mmCreate.CreateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmCreate.CreateMock.defaultExpectation.Counter, 1)
		mm_want := mmCreate.CreateMock.defaultExpectation.params
		mm_got := OrdersRepositoryMockCreateParams{ctx, userID, items}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmCreate.t.Errorf("OrdersRepositoryMock.Create got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmCreate.CreateMock.defaultExpectation.results
		if mm_results == nil {
			mmCreate.t.Fatal("No results are set for the OrdersRepositoryMock.Create")
		}
		return (*mm_results).op1, (*mm_results).err
	}
	if mmCreate.funcCreate != nil {
		return mmCreate.funcCreate(ctx, userID, items)
	}
	mmCreate.t.Fatalf("Unexpected call to OrdersRepositoryMock.Create. %v %v %v", ctx, userID, items)
	return
}

// CreateAfterCounter returns a count of finished OrdersRepositoryMock.Create invocations
func (mmCreate *OrdersRepositoryMock) CreateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCreate.afterCreateCounter)
}

// CreateBeforeCounter returns a count of OrdersRepositoryMock.Create invocations
func (mmCreate *OrdersRepositoryMock) CreateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCreate.beforeCreateCounter)
}

// Calls returns a list of arguments used in each call to OrdersRepositoryMock.Create.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmCreate *mOrdersRepositoryMockCreate) Calls() []*OrdersRepositoryMockCreateParams {
	mmCreate.mutex.RLock()

	argCopy := make([]*OrdersRepositoryMockCreateParams, len(mmCreate.callArgs))
	copy(argCopy, mmCreate.callArgs)

	mmCreate.mutex.RUnlock()

	return argCopy
}

// MinimockCreateDone returns true if the count of the Create invocations corresponds
// the number of defined expectations
func (m *OrdersRepositoryMock) MinimockCreateDone() bool {
	for _, e := range m.CreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CreateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCreateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCreate != nil && mm_atomic.LoadUint64(&m.afterCreateCounter) < 1 {
		return false
	}
	return true
}

// MinimockCreateInspect logs each unmet expectation
func (m *OrdersRepositoryMock) MinimockCreateInspect() {
	for _, e := range m.CreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to OrdersRepositoryMock.Create with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CreateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCreateCounter) < 1 {
		if m.CreateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to OrdersRepositoryMock.Create")
		} else {
			m.t.Errorf("Expected call to OrdersRepositoryMock.Create with params: %#v", *m.CreateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCreate != nil && mm_atomic.LoadUint64(&m.afterCreateCounter) < 1 {
		m.t.Error("Expected call to OrdersRepositoryMock.Create")
	}
}

type mOrdersRepositoryMockGetByID struct {
	mock               *OrdersRepositoryMock
	defaultExpectation *OrdersRepositoryMockGetByIDExpectation
	expectations       []*OrdersRepositoryMockGetByIDExpectation

	callArgs []*OrdersRepositoryMockGetByIDParams
	mutex    sync.RWMutex
}

// OrdersRepositoryMockGetByIDExpectation specifies expectation struct of the OrdersRepository.GetByID
type OrdersRepositoryMockGetByIDExpectation struct {
	mock    *OrdersRepositoryMock
	params  *OrdersRepositoryMockGetByIDParams
	results *OrdersRepositoryMockGetByIDResults
	Counter uint64
}

// OrdersRepositoryMockGetByIDParams contains parameters of the OrdersRepository.GetByID
type OrdersRepositoryMockGetByIDParams struct {
	ctx context.Context
	id  model.OrderID
}

// OrdersRepositoryMockGetByIDResults contains results of the OrdersRepository.GetByID
type OrdersRepositoryMockGetByIDResults struct {
	op1 *model.Order
	err error
}

// Expect sets up expected params for OrdersRepository.GetByID
func (mmGetByID *mOrdersRepositoryMockGetByID) Expect(ctx context.Context, id model.OrderID) *mOrdersRepositoryMockGetByID {
	if mmGetByID.mock.funcGetByID != nil {
		mmGetByID.mock.t.Fatalf("OrdersRepositoryMock.GetByID mock is already set by Set")
	}

	if mmGetByID.defaultExpectation == nil {
		mmGetByID.defaultExpectation = &OrdersRepositoryMockGetByIDExpectation{}
	}

	mmGetByID.defaultExpectation.params = &OrdersRepositoryMockGetByIDParams{ctx, id}
	for _, e := range mmGetByID.expectations {
		if minimock.Equal(e.params, mmGetByID.defaultExpectation.params) {
			mmGetByID.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetByID.defaultExpectation.params)
		}
	}

	return mmGetByID
}

// Inspect accepts an inspector function that has same arguments as the OrdersRepository.GetByID
func (mmGetByID *mOrdersRepositoryMockGetByID) Inspect(f func(ctx context.Context, id model.OrderID)) *mOrdersRepositoryMockGetByID {
	if mmGetByID.mock.inspectFuncGetByID != nil {
		mmGetByID.mock.t.Fatalf("Inspect function is already set for OrdersRepositoryMock.GetByID")
	}

	mmGetByID.mock.inspectFuncGetByID = f

	return mmGetByID
}

// Return sets up results that will be returned by OrdersRepository.GetByID
func (mmGetByID *mOrdersRepositoryMockGetByID) Return(op1 *model.Order, err error) *OrdersRepositoryMock {
	if mmGetByID.mock.funcGetByID != nil {
		mmGetByID.mock.t.Fatalf("OrdersRepositoryMock.GetByID mock is already set by Set")
	}

	if mmGetByID.defaultExpectation == nil {
		mmGetByID.defaultExpectation = &OrdersRepositoryMockGetByIDExpectation{mock: mmGetByID.mock}
	}
	mmGetByID.defaultExpectation.results = &OrdersRepositoryMockGetByIDResults{op1, err}
	return mmGetByID.mock
}

// Set uses given function f to mock the OrdersRepository.GetByID method
func (mmGetByID *mOrdersRepositoryMockGetByID) Set(f func(ctx context.Context, id model.OrderID) (op1 *model.Order, err error)) *OrdersRepositoryMock {
	if mmGetByID.defaultExpectation != nil {
		mmGetByID.mock.t.Fatalf("Default expectation is already set for the OrdersRepository.GetByID method")
	}

	if len(mmGetByID.expectations) > 0 {
		mmGetByID.mock.t.Fatalf("Some expectations are already set for the OrdersRepository.GetByID method")
	}

	mmGetByID.mock.funcGetByID = f
	return mmGetByID.mock
}

// When sets expectation for the OrdersRepository.GetByID which will trigger the result defined by the following
// Then helper
func (mmGetByID *mOrdersRepositoryMockGetByID) When(ctx context.Context, id model.OrderID) *OrdersRepositoryMockGetByIDExpectation {
	if mmGetByID.mock.funcGetByID != nil {
		mmGetByID.mock.t.Fatalf("OrdersRepositoryMock.GetByID mock is already set by Set")
	}

	expectation := &OrdersRepositoryMockGetByIDExpectation{
		mock:   mmGetByID.mock,
		params: &OrdersRepositoryMockGetByIDParams{ctx, id},
	}
	mmGetByID.expectations = append(mmGetByID.expectations, expectation)
	return expectation
}

// Then sets up OrdersRepository.GetByID return parameters for the expectation previously defined by the When method
func (e *OrdersRepositoryMockGetByIDExpectation) Then(op1 *model.Order, err error) *OrdersRepositoryMock {
	e.results = &OrdersRepositoryMockGetByIDResults{op1, err}
	return e.mock
}

// GetByID implements order.OrdersRepository
func (mmGetByID *OrdersRepositoryMock) GetByID(ctx context.Context, id model.OrderID) (op1 *model.Order, err error) {
	mm_atomic.AddUint64(&mmGetByID.beforeGetByIDCounter, 1)
	defer mm_atomic.AddUint64(&mmGetByID.afterGetByIDCounter, 1)

	if mmGetByID.inspectFuncGetByID != nil {
		mmGetByID.inspectFuncGetByID(ctx, id)
	}

	mm_params := OrdersRepositoryMockGetByIDParams{ctx, id}

	// Record call args
	mmGetByID.GetByIDMock.mutex.Lock()
	mmGetByID.GetByIDMock.callArgs = append(mmGetByID.GetByIDMock.callArgs, &mm_params)
	mmGetByID.GetByIDMock.mutex.Unlock()

	for _, e := range mmGetByID.GetByIDMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.op1, e.results.err
		}
	}

	if mmGetByID.GetByIDMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetByID.GetByIDMock.defaultExpectation.Counter, 1)
		mm_want := mmGetByID.GetByIDMock.defaultExpectation.params
		mm_got := OrdersRepositoryMockGetByIDParams{ctx, id}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetByID.t.Errorf("OrdersRepositoryMock.GetByID got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetByID.GetByIDMock.defaultExpectation.results
		if mm_results == nil {
			mmGetByID.t.Fatal("No results are set for the OrdersRepositoryMock.GetByID")
		}
		return (*mm_results).op1, (*mm_results).err
	}
	if mmGetByID.funcGetByID != nil {
		return mmGetByID.funcGetByID(ctx, id)
	}
	mmGetByID.t.Fatalf("Unexpected call to OrdersRepositoryMock.GetByID. %v %v", ctx, id)
	return
}

// GetByIDAfterCounter returns a count of finished OrdersRepositoryMock.GetByID invocations
func (mmGetByID *OrdersRepositoryMock) GetByIDAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetByID.afterGetByIDCounter)
}

// GetByIDBeforeCounter returns a count of OrdersRepositoryMock.GetByID invocations
func (mmGetByID *OrdersRepositoryMock) GetByIDBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetByID.beforeGetByIDCounter)
}

// Calls returns a list of arguments used in each call to OrdersRepositoryMock.GetByID.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetByID *mOrdersRepositoryMockGetByID) Calls() []*OrdersRepositoryMockGetByIDParams {
	mmGetByID.mutex.RLock()

	argCopy := make([]*OrdersRepositoryMockGetByIDParams, len(mmGetByID.callArgs))
	copy(argCopy, mmGetByID.callArgs)

	mmGetByID.mutex.RUnlock()

	return argCopy
}

// MinimockGetByIDDone returns true if the count of the GetByID invocations corresponds
// the number of defined expectations
func (m *OrdersRepositoryMock) MinimockGetByIDDone() bool {
	for _, e := range m.GetByIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetByIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetByIDCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetByID != nil && mm_atomic.LoadUint64(&m.afterGetByIDCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetByIDInspect logs each unmet expectation
func (m *OrdersRepositoryMock) MinimockGetByIDInspect() {
	for _, e := range m.GetByIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to OrdersRepositoryMock.GetByID with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetByIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetByIDCounter) < 1 {
		if m.GetByIDMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to OrdersRepositoryMock.GetByID")
		} else {
			m.t.Errorf("Expected call to OrdersRepositoryMock.GetByID with params: %#v", *m.GetByIDMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetByID != nil && mm_atomic.LoadUint64(&m.afterGetByIDCounter) < 1 {
		m.t.Error("Expected call to OrdersRepositoryMock.GetByID")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *OrdersRepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockCreateInspect()

			m.MinimockGetByIDInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *OrdersRepositoryMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *OrdersRepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockCreateDone() &&
		m.MinimockGetByIDDone()
}
