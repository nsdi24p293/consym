// Code generated by counterfeiter. DO NOT EDIT.
package mocks

import (
	"sync"

	validation "github.com/Yunpeng-J/HLF-2.2/core/handlers/validation/api/state"
)

type StateFetcher struct {
	FetchStateStub        func() (validation.State, error)
	fetchStateMutex       sync.RWMutex
	fetchStateArgsForCall []struct {
	}
	fetchStateReturns struct {
		result1 validation.State
		result2 error
	}
	fetchStateReturnsOnCall map[int]struct {
		result1 validation.State
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *StateFetcher) FetchState() (validation.State, error) {
	fake.fetchStateMutex.Lock()
	ret, specificReturn := fake.fetchStateReturnsOnCall[len(fake.fetchStateArgsForCall)]
	fake.fetchStateArgsForCall = append(fake.fetchStateArgsForCall, struct {
	}{})
	fake.recordInvocation("FetchState", []interface{}{})
	fake.fetchStateMutex.Unlock()
	if fake.FetchStateStub != nil {
		return fake.FetchStateStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.fetchStateReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *StateFetcher) FetchStateCallCount() int {
	fake.fetchStateMutex.RLock()
	defer fake.fetchStateMutex.RUnlock()
	return len(fake.fetchStateArgsForCall)
}

func (fake *StateFetcher) FetchStateCalls(stub func() (validation.State, error)) {
	fake.fetchStateMutex.Lock()
	defer fake.fetchStateMutex.Unlock()
	fake.FetchStateStub = stub
}

func (fake *StateFetcher) FetchStateReturns(result1 validation.State, result2 error) {
	fake.fetchStateMutex.Lock()
	defer fake.fetchStateMutex.Unlock()
	fake.FetchStateStub = nil
	fake.fetchStateReturns = struct {
		result1 validation.State
		result2 error
	}{result1, result2}
}

func (fake *StateFetcher) FetchStateReturnsOnCall(i int, result1 validation.State, result2 error) {
	fake.fetchStateMutex.Lock()
	defer fake.fetchStateMutex.Unlock()
	fake.FetchStateStub = nil
	if fake.fetchStateReturnsOnCall == nil {
		fake.fetchStateReturnsOnCall = make(map[int]struct {
			result1 validation.State
			result2 error
		})
	}
	fake.fetchStateReturnsOnCall[i] = struct {
		result1 validation.State
		result2 error
	}{result1, result2}
}

func (fake *StateFetcher) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.fetchStateMutex.RLock()
	defer fake.fetchStateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *StateFetcher) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
