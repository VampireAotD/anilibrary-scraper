package event

import (
	"context"
	"testing"
	"time"

	"anilibrary-scraper/internal/domain/repository"
	"anilibrary-scraper/internal/domain/repository/models"
	"anilibrary-scraper/internal/domain/service"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type EventServiceSuite struct {
	suite.Suite

	repositoryMock *repository.MockEventRepositoryMockRecorder
	service        service.EventService
}

func TestEventServiceSuite(t *testing.T) {
	suite.Run(t, new(EventServiceSuite))
}

func (suite *EventServiceSuite) SetupSuite() {
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	repositoryMock := repository.NewMockEventRepository(ctrl)

	suite.repositoryMock = repositoryMock.EXPECT()
	suite.service = NewService(repositoryMock)
}

func (suite *EventServiceSuite) TestSend() {
	var (
		t       = suite.T()
		require = suite.Require()
	)

	t.Run("Send message", func(t *testing.T) {
		const testURL string = "https://google.com/"

		suite.repositoryMock.Send(gomock.Any(), models.Event{
			URL:  testURL,
			Date: time.Now().Unix(),
		}).Return(nil)

		err := suite.service.Send(context.Background(), testURL)
		require.NoError(err)
	})
}
