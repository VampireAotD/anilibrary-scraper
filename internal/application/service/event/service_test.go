package event

import (
	"context"
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

func (es *EventServiceSuite) SetupSuite() {
	ctrl := gomock.NewController(es.T())
	defer ctrl.Finish()

	repositoryMock := NewMockRepository(ctrl)

	es.repositoryMock = repositoryMock.EXPECT()
	es.service = NewService(repositoryMock)
}

func (es *EventServiceSuite) TestSend() {
	var (
		t       = es.T()
		require = es.Require()
	)

	t.Run("Send message", func(_ *testing.T) {
		expected := DTO{
			URL:       "https://google.com",
			Time:      time.Now(),
			IP:        "0.0.0.0",
			UserAgent: "Mozilla/5.0",
		}

		es.repositoryMock.Send(gomock.Any(), model.Event{
			URL:       expected.URL,
			Timestamp: expected.Time.Unix(),
			IP:        expected.IP,
			UserAgent: expected.UserAgent,
		}).Return(nil)

		err := es.service.Send(context.Background(), expected)
		require.NoError(err)
	})
}
