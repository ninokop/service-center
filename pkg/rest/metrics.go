//Copyright 2017 Huawei Technologies Co., Ltd
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
package rest

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	handlerName = "servicecenter"
)

var (
	incomingRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "request_total",
			Help:      "Counter of requests received into ROA handler",
		}, []string{"method", "code"})
)

func init() {
	prometheus.MustRegister(incomingRequests)
}

func reportRequestReceived(method string, header http.Header) {
	incomingRequests.WithLabelValues(method, codeFromHeader(header)).Inc()
}

func codeFromHeader(h http.Header) string {
	if code := h.Get("X-Response-Status"); code != "" {
		return code
	}
	return "200"
}
