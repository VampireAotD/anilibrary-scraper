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
	var require = s.Require()

	s.Run("HTTP status", func() {
		s.Run("OK", func() {
			response, err := s.client.fetch(context.Background(), s.mockServer.URL)
			require.NoError(err)
			require.NotNil(response)
		})

		s.Run("404", func() {
			response, err := s.client.fetch(context.Background(), s.mockServer.URL+"/404")
			require.Error(err)
			require.Nil(response)
		})
	})

	s.Run("Concurrency", func() {
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
	s.Run("Parsed", func() {
		doc, err := s.client.HTML(context.Background(), s.mockServer.URL+"/html")
		s.Require().NoError(err)
		s.Require().NotNil(doc)
	})

	s.Run("Wrong content type", func() {
		doc, err := s.client.HTML(context.Background(), s.mockServer.URL+"/404")
		s.Require().Error(err)
		s.Require().Nil(doc)
	})
}
