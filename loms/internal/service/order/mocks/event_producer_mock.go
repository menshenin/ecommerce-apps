// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package mocks

//go:generate minimock -i route256.ozon.ru/project/loms/internal/service/order.EventProducer -o event_producer_mock.go -n EventProducerMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/loms/internal/event"
)

// EventProducerMock implements order.EventProducer
type EventProducerMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcProduce          func(ctx context.Context, event event.OrderEvent)
	inspectFuncProduce   func(ctx context.Context, event event.OrderEvent)
	afterProduceCounter  uint64
	beforeProduceCounter uint64
	ProduceMock          mEventProducerMockProduce
}

// NewEventProducerMock returns a mock for order.EventProducer
func NewEventProducerMock(t minimock.Tester) *EventProducerMock {
	m := &EventProducerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ProduceMock = mEventProducerMockProduce{mock: m}
	m.ProduceMock.callArgs = []*EventProducerMockProduceParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mEventProducerMockProduce struct {
	mock               *EventProducerMock
	defaultExpectation *EventProducerMockProduceExpectation
	expectations       []*EventProducerMockProduceExpectation

	callArgs []*EventProducerMockProduceParams
	mutex    sync.RWMutex
}

// EventProducerMockProduceExpectation specifies expectation struct of the EventProducer.Produce
type EventProducerMockProduceExpectation struct {
	mock   *EventProducerMock
	params *EventProducerMockProduceParams

	Counter uint64
}

// EventProducerMockProduceParams contains parameters of the EventProducer.Produce
type EventProducerMockProduceParams struct {
	ctx   context.Context
	event event.OrderEvent
}

// Expect sets up expected params for EventProducer.Produce
func (mmProduce *mEventProducerMockProduce) Expect(ctx context.Context, event event.OrderEvent) *mEventProducerMockProduce {
	if mmProduce.mock.funcProduce != nil {
		mmProduce.mock.t.Fatalf("EventProducerMock.Produce mock is already set by Set")
	}

	if mmProduce.defaultExpectation == nil {
		mmProduce.defaultExpectation = &EventProducerMockProduceExpectation{}
	}

	mmProduce.defaultExpectation.params = &EventProducerMockProduceParams{ctx, event}
	for _, e := range mmProduce.expectations {
		if minimock.Equal(e.params, mmProduce.defaultExpectation.params) {
			mmProduce.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmProduce.defaultExpectation.params)
		}
	}

	return mmProduce
}

// Inspect accepts an inspector function that has same arguments as the EventProducer.Produce
func (mmProduce *mEventProducerMockProduce) Inspect(f func(ctx context.Context, event event.OrderEvent)) *mEventProducerMockProduce {
	if mmProduce.mock.inspectFuncProduce != nil {
		mmProduce.mock.t.Fatalf("Inspect function is already set for EventProducerMock.Produce")
	}

	mmProduce.mock.inspectFuncProduce = f

	return mmProduce
}

// Return sets up results that will be returned by EventProducer.Produce
func (mmProduce *mEventProducerMockProduce) Return() *EventProducerMock {
	if mmProduce.mock.funcProduce != nil {
		mmProduce.mock.t.Fatalf("EventProducerMock.Produce mock is already set by Set")
	}

	if mmProduce.defaultExpectation == nil {
		mmProduce.defaultExpectation = &EventProducerMockProduceExpectation{mock: mmProduce.mock}
	}

	return mmProduce.mock
}

// Set uses given function f to mock the EventProducer.Produce method
func (mmProduce *mEventProducerMockProduce) Set(f func(ctx context.Context, event event.OrderEvent)) *EventProducerMock {
	if mmProduce.defaultExpectation != nil {
		mmProduce.mock.t.Fatalf("Default expectation is already set for the EventProducer.Produce method")
	}

	if len(mmProduce.expectations) > 0 {
		mmProduce.mock.t.Fatalf("Some expectations are already set for the EventProducer.Produce method")
	}

	mmProduce.mock.funcProduce = f
	return mmProduce.mock
}

// Produce implements order.EventProducer
func (mmProduce *EventProducerMock) Produce(ctx context.Context, event event.OrderEvent) {
	mm_atomic.AddUint64(&mmProduce.beforeProduceCounter, 1)
	defer mm_atomic.AddUint64(&mmProduce.afterProduceCounter, 1)

	if mmProduce.inspectFuncProduce != nil {
		mmProduce.inspectFuncProduce(ctx, event)
	}

	mm_params := EventProducerMockProduceParams{ctx, event}

	// Record call args
	mmProduce.ProduceMock.mutex.Lock()
	mmProduce.ProduceMock.callArgs = append(mmProduce.ProduceMock.callArgs, &mm_params)
	mmProduce.ProduceMock.mutex.Unlock()

	for _, e := range mmProduce.ProduceMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmProduce.ProduceMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmProduce.ProduceMock.defaultExpectation.Counter, 1)
		mm_want := mmProduce.ProduceMock.defaultExpectation.params
		mm_got := EventProducerMockProduceParams{ctx, event}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmProduce.t.Errorf("EventProducerMock.Produce got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		return

	}
	if mmProduce.funcProduce != nil {
		mmProduce.funcProduce(ctx, event)
		return
	}
	mmProduce.t.Fatalf("Unexpected call to EventProducerMock.Produce. %v %v", ctx, event)

}

// ProduceAfterCounter returns a count of finished EventProducerMock.Produce invocations
func (mmProduce *EventProducerMock) ProduceAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmProduce.afterProduceCounter)
}

// ProduceBeforeCounter returns a count of EventProducerMock.Produce invocations
func (mmProduce *EventProducerMock) ProduceBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmProduce.beforeProduceCounter)
}

// Calls returns a list of arguments used in each call to EventProducerMock.Produce.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmProduce *mEventProducerMockProduce) Calls() []*EventProducerMockProduceParams {
	mmProduce.mutex.RLock()

	argCopy := make([]*EventProducerMockProduceParams, len(mmProduce.callArgs))
	copy(argCopy, mmProduce.callArgs)

	mmProduce.mutex.RUnlock()

	return argCopy
}

// MinimockProduceDone returns true if the count of the Produce invocations corresponds
// the number of defined expectations
func (m *EventProducerMock) MinimockProduceDone() bool {
	for _, e := range m.ProduceMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ProduceMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterProduceCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcProduce != nil && mm_atomic.LoadUint64(&m.afterProduceCounter) < 1 {
		return false
	}
	return true
}

// MinimockProduceInspect logs each unmet expectation
func (m *EventProducerMock) MinimockProduceInspect() {
	for _, e := range m.ProduceMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to EventProducerMock.Produce with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ProduceMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterProduceCounter) < 1 {
		if m.ProduceMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to EventProducerMock.Produce")
		} else {
			m.t.Errorf("Expected call to EventProducerMock.Produce with params: %#v", *m.ProduceMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcProduce != nil && mm_atomic.LoadUint64(&m.afterProduceCounter) < 1 {
		m.t.Error("Expected call to EventProducerMock.Produce")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *EventProducerMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockProduceInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *EventProducerMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *EventProducerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockProduceDone()
}
