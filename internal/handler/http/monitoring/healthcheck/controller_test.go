package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
)

type HealthcheckControllerSuite struct {
	suite.Suite

	redisServer *miniredis.Miniredis
	controller  Controller
}

func TestHealthcheckControllerSuite(t *testing.T) {
	suite.Run(t, new(HealthcheckControllerSuite))
}

func (suite *HealthcheckControllerSuite) SetupSuite() {
	suite.redisServer = miniredis.RunT(suite.T())
	redisClient := redis.NewClient(&redis.Options{
		Addr: suite.redisServer.Addr(),
	})

	suite.controller = NewController(redisClient)
}

func (suite *HealthcheckControllerSuite) TearDownSuite() {
	suite.redisServer.Close()
}

func (suite *HealthcheckControllerSuite) sendHealthcheckRequest() *httptest.ResponseRecorder {
	handler := suite.controller.Healthcheck
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)

	handler(recorder, request)

	return recorder
}

func (suite *HealthcheckControllerSuite) TestHealthcheck() {
	t := suite.T()
	require := suite.Require()

	t.Run("Redis up", func(t *testing.T) {
		response := suite.sendHealthcheckRequest()
		require.Equal(http.StatusOK, response.Code)
	})

	t.Run("Redis down", func(t *testing.T) {
		suite.redisServer.Close()

		response := suite.sendHealthcheckRequest()
		require.Equal(http.StatusInternalServerError, response.Code)
	})
}
