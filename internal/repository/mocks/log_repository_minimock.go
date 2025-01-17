// Code generated by http://github.com/gojuno/minimock (v3.4.1). DO NOT EDIT.

package mocks

//go:generate minimock -i github.com/Mobo140/chat/internal/repository.LogRepository -o log_repository_minimock.go -n LogRepositoryMock -p mocks

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/Mobo140/chat/internal/model"
	"github.com/gojuno/minimock/v3"
)

// LogRepositoryMock implements mm_repository.LogRepository
type LogRepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcCreate          func(ctx context.Context, logEntry *model.LogEntry) (err error)
	funcCreateOrigin    string
	inspectFuncCreate   func(ctx context.Context, logEntry *model.LogEntry)
	afterCreateCounter  uint64
	beforeCreateCounter uint64
	CreateMock          mLogRepositoryMockCreate
}

// NewLogRepositoryMock returns a mock for mm_repository.LogRepository
func NewLogRepositoryMock(t minimock.Tester) *LogRepositoryMock {
	m := &LogRepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.CreateMock = mLogRepositoryMockCreate{mock: m}
	m.CreateMock.callArgs = []*LogRepositoryMockCreateParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mLogRepositoryMockCreate struct {
	optional           bool
	mock               *LogRepositoryMock
	defaultExpectation *LogRepositoryMockCreateExpectation
	expectations       []*LogRepositoryMockCreateExpectation

	callArgs []*LogRepositoryMockCreateParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// LogRepositoryMockCreateExpectation specifies expectation struct of the LogRepository.Create
type LogRepositoryMockCreateExpectation struct {
	mock               *LogRepositoryMock
	params             *LogRepositoryMockCreateParams
	paramPtrs          *LogRepositoryMockCreateParamPtrs
	expectationOrigins LogRepositoryMockCreateExpectationOrigins
	results            *LogRepositoryMockCreateResults
	returnOrigin       string
	Counter            uint64
}

// LogRepositoryMockCreateParams contains parameters of the LogRepository.Create
type LogRepositoryMockCreateParams struct {
	ctx      context.Context
	logEntry *model.LogEntry
}

// LogRepositoryMockCreateParamPtrs contains pointers to parameters of the LogRepository.Create
type LogRepositoryMockCreateParamPtrs struct {
	ctx      *context.Context
	logEntry **model.LogEntry
}

// LogRepositoryMockCreateResults contains results of the LogRepository.Create
type LogRepositoryMockCreateResults struct {
	err error
}

// LogRepositoryMockCreateOrigins contains origins of expectations of the LogRepository.Create
type LogRepositoryMockCreateExpectationOrigins struct {
	origin         string
	originCtx      string
	originLogEntry string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmCreate *mLogRepositoryMockCreate) Optional() *mLogRepositoryMockCreate {
	mmCreate.optional = true
	return mmCreate
}

// Expect sets up expected params for LogRepository.Create
func (mmCreate *mLogRepositoryMockCreate) Expect(ctx context.Context, logEntry *model.LogEntry) *mLogRepositoryMockCreate {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("LogRepositoryMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &LogRepositoryMockCreateExpectation{}
	}

	if mmCreate.defaultExpectation.paramPtrs != nil {
		mmCreate.mock.t.Fatalf("LogRepositoryMock.Create mock is already set by ExpectParams functions")
	}

	mmCreate.defaultExpectation.params = &LogRepositoryMockCreateParams{ctx, logEntry}
	mmCreate.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmCreate.expectations {
		if minimock.Equal(e.params, mmCreate.defaultExpectation.params) {
			mmCreate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmCreate.defaultExpectation.params)
		}
	}

	return mmCreate
}

// ExpectCtxParam1 sets up expected param ctx for LogRepository.Create
func (mmCreate *mLogRepositoryMockCreate) ExpectCtxParam1(ctx context.Context) *mLogRepositoryMockCreate {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("LogRepositoryMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &LogRepositoryMockCreateExpectation{}
	}

	if mmCreate.defaultExpectation.params != nil {
		mmCreate.mock.t.Fatalf("LogRepositoryMock.Create mock is already set by Expect")
	}

	if mmCreate.defaultExpectation.paramPtrs == nil {
		mmCreate.defaultExpectation.paramPtrs = &LogRepositoryMockCreateParamPtrs{}
	}
	mmCreate.defaultExpectation.paramPtrs.ctx = &ctx
	mmCreate.defaultExpectation.expectationOrigins.originCtx = minimock.CallerInfo(1)

	return mmCreate
}

// ExpectLogEntryParam2 sets up expected param logEntry for LogRepository.Create
func (mmCreate *mLogRepositoryMockCreate) ExpectLogEntryParam2(logEntry *model.LogEntry) *mLogRepositoryMockCreate {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("LogRepositoryMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &LogRepositoryMockCreateExpectation{}
	}

	if mmCreate.defaultExpectation.params != nil {
		mmCreate.mock.t.Fatalf("LogRepositoryMock.Create mock is already set by Expect")
	}

	if mmCreate.defaultExpectation.paramPtrs == nil {
		mmCreate.defaultExpectation.paramPtrs = &LogRepositoryMockCreateParamPtrs{}
	}
	mmCreate.defaultExpectation.paramPtrs.logEntry = &logEntry
	mmCreate.defaultExpectation.expectationOrigins.originLogEntry = minimock.CallerInfo(1)

	return mmCreate
}

// Inspect accepts an inspector function that has same arguments as the LogRepository.Create
func (mmCreate *mLogRepositoryMockCreate) Inspect(f func(ctx context.Context, logEntry *model.LogEntry)) *mLogRepositoryMockCreate {
	if mmCreate.mock.inspectFuncCreate != nil {
		mmCreate.mock.t.Fatalf("Inspect function is already set for LogRepositoryMock.Create")
	}

	mmCreate.mock.inspectFuncCreate = f

	return mmCreate
}

// Return sets up results that will be returned by LogRepository.Create
func (mmCreate *mLogRepositoryMockCreate) Return(err error) *LogRepositoryMock {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("LogRepositoryMock.Create mock is already set by Set")
	}

	if mmCreate.defaultExpectation == nil {
		mmCreate.defaultExpectation = &LogRepositoryMockCreateExpectation{mock: mmCreate.mock}
	}
	mmCreate.defaultExpectation.results = &LogRepositoryMockCreateResults{err}
	mmCreate.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmCreate.mock
}

// Set uses given function f to mock the LogRepository.Create method
func (mmCreate *mLogRepositoryMockCreate) Set(f func(ctx context.Context, logEntry *model.LogEntry) (err error)) *LogRepositoryMock {
	if mmCreate.defaultExpectation != nil {
		mmCreate.mock.t.Fatalf("Default expectation is already set for the LogRepository.Create method")
	}

	if len(mmCreate.expectations) > 0 {
		mmCreate.mock.t.Fatalf("Some expectations are already set for the LogRepository.Create method")
	}

	mmCreate.mock.funcCreate = f
	mmCreate.mock.funcCreateOrigin = minimock.CallerInfo(1)
	return mmCreate.mock
}

// When sets expectation for the LogRepository.Create which will trigger the result defined by the following
// Then helper
func (mmCreate *mLogRepositoryMockCreate) When(ctx context.Context, logEntry *model.LogEntry) *LogRepositoryMockCreateExpectation {
	if mmCreate.mock.funcCreate != nil {
		mmCreate.mock.t.Fatalf("LogRepositoryMock.Create mock is already set by Set")
	}

	expectation := &LogRepositoryMockCreateExpectation{
		mock:               mmCreate.mock,
		params:             &LogRepositoryMockCreateParams{ctx, logEntry},
		expectationOrigins: LogRepositoryMockCreateExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmCreate.expectations = append(mmCreate.expectations, expectation)
	return expectation
}

// Then sets up LogRepository.Create return parameters for the expectation previously defined by the When method
func (e *LogRepositoryMockCreateExpectation) Then(err error) *LogRepositoryMock {
	e.results = &LogRepositoryMockCreateResults{err}
	return e.mock
}

// Times sets number of times LogRepository.Create should be invoked
func (mmCreate *mLogRepositoryMockCreate) Times(n uint64) *mLogRepositoryMockCreate {
	if n == 0 {
		mmCreate.mock.t.Fatalf("Times of LogRepositoryMock.Create mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmCreate.expectedInvocations, n)
	mmCreate.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmCreate
}

func (mmCreate *mLogRepositoryMockCreate) invocationsDone() bool {
	if len(mmCreate.expectations) == 0 && mmCreate.defaultExpectation == nil && mmCreate.mock.funcCreate == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmCreate.mock.afterCreateCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmCreate.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Create implements mm_repository.LogRepository
func (mmCreate *LogRepositoryMock) Create(ctx context.Context, logEntry *model.LogEntry) (err error) {
	mm_atomic.AddUint64(&mmCreate.beforeCreateCounter, 1)
	defer mm_atomic.AddUint64(&mmCreate.afterCreateCounter, 1)

	mmCreate.t.Helper()

	if mmCreate.inspectFuncCreate != nil {
		mmCreate.inspectFuncCreate(ctx, logEntry)
	}

	mm_params := LogRepositoryMockCreateParams{ctx, logEntry}

	// Record call args
	mmCreate.CreateMock.mutex.Lock()
	mmCreate.CreateMock.callArgs = append(mmCreate.CreateMock.callArgs, &mm_params)
	mmCreate.CreateMock.mutex.Unlock()

	for _, e := range mmCreate.CreateMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmCreate.CreateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmCreate.CreateMock.defaultExpectation.Counter, 1)
		mm_want := mmCreate.CreateMock.defaultExpectation.params
		mm_want_ptrs := mmCreate.CreateMock.defaultExpectation.paramPtrs

		mm_got := LogRepositoryMockCreateParams{ctx, logEntry}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.ctx != nil && !minimock.Equal(*mm_want_ptrs.ctx, mm_got.ctx) {
				mmCreate.t.Errorf("LogRepositoryMock.Create got unexpected parameter ctx, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmCreate.CreateMock.defaultExpectation.expectationOrigins.originCtx, *mm_want_ptrs.ctx, mm_got.ctx, minimock.Diff(*mm_want_ptrs.ctx, mm_got.ctx))
			}

			if mm_want_ptrs.logEntry != nil && !minimock.Equal(*mm_want_ptrs.logEntry, mm_got.logEntry) {
				mmCreate.t.Errorf("LogRepositoryMock.Create got unexpected parameter logEntry, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmCreate.CreateMock.defaultExpectation.expectationOrigins.originLogEntry, *mm_want_ptrs.logEntry, mm_got.logEntry, minimock.Diff(*mm_want_ptrs.logEntry, mm_got.logEntry))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmCreate.t.Errorf("LogRepositoryMock.Create got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmCreate.CreateMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmCreate.CreateMock.defaultExpectation.results
		if mm_results == nil {
			mmCreate.t.Fatal("No results are set for the LogRepositoryMock.Create")
		}
		return (*mm_results).err
	}
	if mmCreate.funcCreate != nil {
		return mmCreate.funcCreate(ctx, logEntry)
	}
	mmCreate.t.Fatalf("Unexpected call to LogRepositoryMock.Create. %v %v", ctx, logEntry)
	return
}

// CreateAfterCounter returns a count of finished LogRepositoryMock.Create invocations
func (mmCreate *LogRepositoryMock) CreateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCreate.afterCreateCounter)
}

// CreateBeforeCounter returns a count of LogRepositoryMock.Create invocations
func (mmCreate *LogRepositoryMock) CreateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCreate.beforeCreateCounter)
}

// Calls returns a list of arguments used in each call to LogRepositoryMock.Create.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmCreate *mLogRepositoryMockCreate) Calls() []*LogRepositoryMockCreateParams {
	mmCreate.mutex.RLock()

	argCopy := make([]*LogRepositoryMockCreateParams, len(mmCreate.callArgs))
	copy(argCopy, mmCreate.callArgs)

	mmCreate.mutex.RUnlock()

	return argCopy
}

// MinimockCreateDone returns true if the count of the Create invocations corresponds
// the number of defined expectations
func (m *LogRepositoryMock) MinimockCreateDone() bool {
	if m.CreateMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.CreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.CreateMock.invocationsDone()
}

// MinimockCreateInspect logs each unmet expectation
func (m *LogRepositoryMock) MinimockCreateInspect() {
	for _, e := range m.CreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LogRepositoryMock.Create at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterCreateCounter := mm_atomic.LoadUint64(&m.afterCreateCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.CreateMock.defaultExpectation != nil && afterCreateCounter < 1 {
		if m.CreateMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to LogRepositoryMock.Create at\n%s", m.CreateMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to LogRepositoryMock.Create at\n%s with params: %#v", m.CreateMock.defaultExpectation.expectationOrigins.origin, *m.CreateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCreate != nil && afterCreateCounter < 1 {
		m.t.Errorf("Expected call to LogRepositoryMock.Create at\n%s", m.funcCreateOrigin)
	}

	if !m.CreateMock.invocationsDone() && afterCreateCounter > 0 {
		m.t.Errorf("Expected %d calls to LogRepositoryMock.Create at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.CreateMock.expectedInvocations), m.CreateMock.expectedInvocationsOrigin, afterCreateCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *LogRepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockCreateInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *LogRepositoryMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *LogRepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockCreateDone()
}
