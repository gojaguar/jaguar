package server

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestServer_TriggerError(t *testing.T) {
	var builder Builder
	srv := builder.Build()

	// Closing the server before ListenAndServe is called should make ListenAndServe return an error
	assert.NoError(t, srv.http.Close())

	err := srv.ListenAndServe()
	assert.Error(t, err)
	assert.EqualError(t, err, "http: Server closed")
}

func TestServer_TriggerSignal(t *testing.T) {
	var builder Builder
	var ts testSignal
	srv := builder.Signal(&ts).Build()

	var wg sync.WaitGroup
	wg.Add(1)

	var err error
	go func(t *testing.T, srv *Server) {
		err = srv.ListenAndServe()
		assert.Error(t, err)
		assert.EqualError(t, err, "signal testSignal triggered")
		wg.Done()
	}(t, srv)

	srv.sigs <- &ts

	wg.Wait()
}

type testSignal struct {
	T *testing.T
}

func (t *testSignal) String() string {
	return "testSignal"
}

func (t *testSignal) Signal() {}
