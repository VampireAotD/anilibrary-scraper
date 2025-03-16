package event

import (
	"testing"
	"time"

	"github.com/VampireAotD/anilibrary-scraper/internal/infrastructure/repository/model"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type EventServiceSuite struct {
	suite.Suite

	repositoryMock *MockRepositoryMockRecorder
	service        Service
}

func TestEventServiceSuite(t *testing.T) {
	suite.Run(t, new(EventServiceSuite))
}

func (s *EventServiceSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	repositoryMock := NewMockRepository(ctrl)

	s.repositoryMock = repositoryMock.EXPECT()
	s.service = NewService(repositoryMock)
}

func (s *EventServiceSuite) TestSend() {
	var require = s.Require()

	s.T().Run("Send message", func(t *testing.T) {
		expected := DTO{
			URL:       "https://google.com",
			Time:      time.Now(),
			IP:        "0.0.0.0",
			UserAgent: "Mozilla/5.0",
		}

		s.repositoryMock.Send(gomock.Any(), model.Event{
			URL:       expected.URL,
			Timestamp: expected.Time.Unix(),
			IP:        expected.IP,
			UserAgent: expected.UserAgent,
		}).Return(nil)

		err := s.service.Send(t.Context(), expected)
		require.NoError(err)
	})
}
