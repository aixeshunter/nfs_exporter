// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/log"
)

const (
	ExecPath = "/usr/sbin/showmount"
)

type nfs struct {
	*httptest.Server
	response []byte
}

func handler(h *nfs) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(h.response)
	}
}

func newNFS(response []byte) *nfs {
	h := &nfs{response: response}
	h.Server = httptest.NewServer(handler(h))
	return h
}

func readGauge(m prometheus.Gauge) float64 {
	// TODO: Revisit this once client_golang offers better testing tools.
	pb := &dto.Metric{}
	m.Write(pb)
	return pb.GetGauge().GetValue()
}

type testCase struct {
	mountpath string
	address   string
}

var tests = []testCase{
	{
		mountpath: "/opt/nfs1, /opt/nfs1, /opt/nfs2",
		address:   "127.0.0.1",
	},
	{
		mountpath: "/opt/nfs1, /opt/nfs1, /opt/nfs2",
		address:   "192.168.0.1",
	},
}

func TestInvalidConfig(t *testing.T) {
	h := newNFS([]byte("not,enough,fields"))
	defer h.Close()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("While trying to get Hostname error happened: %v", err)
	}

	e, _ := NewExporter(hostname, ExecPath, tests[1].mountpath, tests[1].address)
	ch := make(chan prometheus.Metric)

	go func() {
		defer close(ch)
		e.Collect(ch)
	}()

	if expect, got := 1., readGauge((<-ch).(prometheus.Gauge)); expect != got {
		// up
		t.Errorf("expected %f up, got %f", expect, got)
	}

	if <-ch != nil {
		t.Errorf("expected closed channel")
	}
}
