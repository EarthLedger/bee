// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ethersphere/bee/pkg/jsonhttp"
	"github.com/ethersphere/bee/pkg/jsonhttp/jsonhttptest"
	"github.com/ethersphere/bee/pkg/logging"
	statestore "github.com/ethersphere/bee/pkg/statestore/mock"
	smock "github.com/ethersphere/bee/pkg/storage/mock"
	"github.com/ethersphere/bee/pkg/swarm"
	"github.com/ethersphere/bee/pkg/tags"
)

func TestStewardshipReUpload(t *testing.T) {
	var (
		logger         = logging.New(ioutil.Discard, 0)
		mockStatestore = statestore.NewStateStore()
		m              = &mockSteward{}
		storer         = smock.NewStorer()
		addr           = swarm.NewAddress([]byte{31: 128})
	)
	client, _, _ := newTestServer(t, testServerOptions{
		Storer:  storer,
		Tags:    tags.NewTags(mockStatestore, logger),
		Logger:  logger,
		Steward: m,
	})
	jsonhttptest.Request(t, client, http.MethodPut, "/v1/stewardship/"+addr.String(), http.StatusOK,
		jsonhttptest.WithExpectedJSONResponse(jsonhttp.StatusResponse{
			Message: http.StatusText(http.StatusOK),
			Code:    http.StatusOK,
		}),
	)
	if !m.addr.Equal(addr) {
		t.Fatalf("\nhave address: %q\nwant address: %q", m.addr.String(), addr.String())
	}
}

type mockSteward struct {
	addr swarm.Address
}

func (m *mockSteward) Reupload(_ context.Context, addr swarm.Address) error {
	m.addr = addr
	return nil
}
