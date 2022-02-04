// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	chaincode "github.com/Yunpeng-J/HLF-2.2/common/chaincode"
	mock "github.com/stretchr/testify/mock"
)

// Enumerator is an autogenerated mock type for the Enumerator type
type Enumerator struct {
	mock.Mock
}

// Enumerate provides a mock function with given fields:
func (_m *Enumerator) Enumerate() ([]chaincode.InstalledChaincode, error) {
	ret := _m.Called()

	var r0 []chaincode.InstalledChaincode
	if rf, ok := ret.Get(0).(func() []chaincode.InstalledChaincode); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]chaincode.InstalledChaincode)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
