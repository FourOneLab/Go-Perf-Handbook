package http_server

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"testing"
)

// TODO: use net/http/httptest

var (
	httpServer  *http.Server
	httpsServer *http.Server
)

type testHandler struct{}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK\n"))
}

func startHTTPServer() {
	if httpServer != nil {
		return
	}
	httpServer = &http.Server{Handler: &testHandler{}}

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := httpServer.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func startHTTPSServer() {
	if httpsServer != nil {
		return
	}

	httpsServer = &http.Server{Handler: &testHandler{}}
	lis, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := httpsServer.ServeTLS(lis, "server.crt", "server.key"); err != nil {
			panic(err)
		}
	}()
}

func sendRequest(client *http.Client, addr string) {
	resp, err := client.Get(addr)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic("request failed")
	}
	if _, err = io.ReadAll(resp.Body); err != nil {
		panic(err)
	}
	if err = resp.Body.Close(); err != nil {
		panic(err)
	}
}

func BenchmarkHTTP(b *testing.B) {
	startHTTPServer()
	client := &http.Client{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendRequest(client, "http://localhost:8081")
	}
}

func BenchmarkHTTPNoKeepAlive(b *testing.B) {
	startHTTPServer()
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sendRequest(client, "http://localhost:8081")
	}
}

func BenchmarkHTTPSNoKeepAlive(b *testing.B) {
	startHTTPSServer()

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sendRequest(client, "https://127.0.0.1:8443/")
	}
}
