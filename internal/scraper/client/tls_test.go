package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TLSClientSuite struct {
	suite.Suite

	client     TLSClient
	mockServer *httptest.Server
}

func TestTLSClientSuite(t *testing.T) {
	suite.Run(t, new(TLSClientSuite))
}

func (s *TLSClientSuite) SetupTest() {
	client, err := NewTLSClient()
	s.Require().NoError(err)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/404", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	mux.HandleFunc("/html", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, writeErr := w.Write([]byte("<html><body><h1>Hello World</h1></body></html>"))
		s.Require().NoError(writeErr)
	})

	s.client = client
	s.mockServer = httptest.NewServer(mux)
}

func (s *TLSClientSuite) TearDownSuite() {
	s.mockServer.Close()
}

func (s *TLSClientSuite) TestFetch() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("HTTP status", func(t *testing.T) {
		t.Run("OK", func(_ *testing.T) {
			response, err := s.client.fetch(context.Background(), s.mockServer.URL)
			require.NoError(err)
			require.NotNil(response)
		})

		t.Run("404", func(_ *testing.T) {
			response, err := s.client.fetch(context.Background(), s.mockServer.URL+"/404")
			require.Error(err)
			require.Nil(response)
		})
	})

	t.Run("Concurrency", func(_ *testing.T) {
		var wg sync.WaitGroup

		wg.Add(2)

		go func() {
			defer wg.Done()
			response, err := s.client.fetch(context.Background(), s.mockServer.URL)
			require.NoError(err)
			require.NotNil(response)
		}()

		go func() {
			defer wg.Done()
			response, err := s.client.fetch(context.Background(), s.mockServer.URL+"/404")
			require.Error(err)
			require.Nil(response)
		}()

		wg.Wait()
	})
}

func (s *TLSClientSuite) TestHTML() {
	var (
		t       = s.T()
		require = s.Require()
	)

	t.Run("Parsed", func(_ *testing.T) {
		doc, err := s.client.HTML(context.Background(), s.mockServer.URL+"/html")
		require.NoError(err)
		require.NotNil(doc)
	})

	t.Run("Wrong content type", func(_ *testing.T) {
		doc, err := s.client.HTML(context.Background(), s.mockServer.URL+"/404")
		require.Error(err)
		require.Nil(doc)
	})
}
