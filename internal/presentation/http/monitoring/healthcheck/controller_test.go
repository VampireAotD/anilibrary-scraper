package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

const endpoint string = "/healthcheck"

type HealthcheckControllerSuite struct {
	controller Controller
	suite.Suite
	router *fiber.App
}

func TestHealthcheckControllerSuite(t *testing.T) {
	suite.Run(t, new(HealthcheckControllerSuite))
}

func (s *HealthcheckControllerSuite) SetupSuite() {
	s.controller = NewController()
	s.router = fiber.New()

	s.router.Get(endpoint, s.controller.Healthcheck)
}

func (s *HealthcheckControllerSuite) TestHealthcheck() {
	request := httptest.NewRequest(http.MethodGet, endpoint, http.NoBody)

	response, err := s.router.Test(request, -1)
	s.Require().NoError(err)
	s.Require().NoError(response.Body.Close())

	s.Require().Equal(http.StatusOK, response.StatusCode)
}
