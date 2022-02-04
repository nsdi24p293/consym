/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package protoext_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/Yunpeng-J/fabric-protos-go/discovery"
	"github.com/Yunpeng-J/HLF-2.2/discovery/protoext"
	"github.com/stretchr/testify/assert"
)

func TestSignedRequestToRequest(t *testing.T) {
	sr := &discovery.SignedRequest{
		Payload: []byte{0},
	}
	_, err := protoext.SignedRequestToRequest(sr)
	assert.Error(t, err)

	req := &discovery.Request{}
	b, _ := proto.Marshal(req)
	sr.Payload = b
	r, err := protoext.SignedRequestToRequest(sr)
	assert.NoError(t, err)
	assert.NotNil(t, r)
}
