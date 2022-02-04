// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	common "github.com/Yunpeng-J/fabric-protos-go/common"
	mock "github.com/stretchr/testify/mock"

	peer "github.com/Yunpeng-J/fabric-protos-go/peer"
)

// Dispatcher is an autogenerated mock type for the Dispatcher type
type Dispatcher struct {
	mock.Mock
}

// Dispatch provides a mock function with given fields: seq, payload, envBytes, block
func (_m *Dispatcher) Dispatch(seq int, payload *common.Payload, envBytes []byte, block *common.Block) (error, peer.TxValidationCode) {
	ret := _m.Called(seq, payload, envBytes, block)

	var r0 error
	if rf, ok := ret.Get(0).(func(int, *common.Payload, []byte, *common.Block) error); ok {
		r0 = rf(seq, payload, envBytes, block)
	} else {
		r0 = ret.Error(0)
	}

	var r1 peer.TxValidationCode
	if rf, ok := ret.Get(1).(func(int, *common.Payload, []byte, *common.Block) peer.TxValidationCode); ok {
		r1 = rf(seq, payload, envBytes, block)
	} else {
		r1 = ret.Get(1).(peer.TxValidationCode)
	}

	return r0, r1
}
