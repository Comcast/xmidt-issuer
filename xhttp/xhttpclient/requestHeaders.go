// SPDX-FileCopyrightText: 2017 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package xhttpclient

import (
	"net/http"

	"github.com/xmidt-org/themis/xhttp"
)

// RequestHeaders provides a RoundTripper constructor that inserts a constant set of headers
// into each request
type RequestHeaders struct {
	Header http.Header
}

func (rh RequestHeaders) Then(next http.RoundTripper) http.RoundTripper {
	if len(rh.Header) == 0 {
		return next
	}

	header := xhttp.CanonicalizeHeaders(rh.Header)
	return RoundTripperFunc(func(request *http.Request) (*http.Response, error) {
		xhttp.SetHeaders(request.Header, header)
		return next.RoundTrip(request)
	})
}

func (rh RequestHeaders) ThenFunc(next RoundTripperFunc) http.RoundTripper {
	return rh.Then(next)
}
