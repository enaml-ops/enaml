// This file was generated by counterfeiter
package releasejob_experimentfakes

import (
	"sync"

	"github.com/enaml-ops/enaml/releasejob_experiment"
)

type FakeBoshJob struct {
	MetaStub        func() releasejob_experiment.BoshJobMeta
	metaMutex       sync.RWMutex
	metaArgsForCall []struct{}
	metaReturns     struct {
		result1 releasejob_experiment.BoshJobMeta
	}
	StartStub        func() error
	startMutex       sync.RWMutex
	startArgsForCall []struct{}
	startReturns     struct {
		result1 error
	}
	StopStub        func() error
	stopMutex       sync.RWMutex
	stopArgsForCall []struct{}
	stopReturns     struct {
		result1 error
	}
}

func (fake *FakeBoshJob) Meta() releasejob_experiment.BoshJobMeta {
	fake.metaMutex.Lock()
	fake.metaArgsForCall = append(fake.metaArgsForCall, struct{}{})
	fake.metaMutex.Unlock()
	if fake.MetaStub != nil {
		return fake.MetaStub()
	} else {
		return fake.metaReturns.result1
	}
}

func (fake *FakeBoshJob) MetaCallCount() int {
	fake.metaMutex.RLock()
	defer fake.metaMutex.RUnlock()
	return len(fake.metaArgsForCall)
}

func (fake *FakeBoshJob) MetaReturns(result1 releasejob_experiment.BoshJobMeta) {
	fake.MetaStub = nil
	fake.metaReturns = struct {
		result1 releasejob_experiment.BoshJobMeta
	}{result1}
}

func (fake *FakeBoshJob) Start() error {
	fake.startMutex.Lock()
	fake.startArgsForCall = append(fake.startArgsForCall, struct{}{})
	fake.startMutex.Unlock()
	if fake.StartStub != nil {
		return fake.StartStub()
	} else {
		return fake.startReturns.result1
	}
}

func (fake *FakeBoshJob) StartCallCount() int {
	fake.startMutex.RLock()
	defer fake.startMutex.RUnlock()
	return len(fake.startArgsForCall)
}

func (fake *FakeBoshJob) StartReturns(result1 error) {
	fake.StartStub = nil
	fake.startReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBoshJob) Stop() error {
	fake.stopMutex.Lock()
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct{}{})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		return fake.StopStub()
	} else {
		return fake.stopReturns.result1
	}
}

func (fake *FakeBoshJob) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *FakeBoshJob) StopReturns(result1 error) {
	fake.StopStub = nil
	fake.stopReturns = struct {
		result1 error
	}{result1}
}

var _ releasejob_experiment.BoshJob = new(FakeBoshJob)