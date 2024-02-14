package storage

import (
	"fmt"
	"github.com/chemax/url-shorter/config"
	"github.com/chemax/url-shorter/logger"
	mock_storage "github.com/chemax/url-shorter/mocks/newstorage"
	"github.com/chemax/url-shorter/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Конечно, тут стоило бы работать с реальным хранилищем. Но долго, а я в цейтноте.
func TestNewStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cfg, err := config.NewConfig()
	assert.Nil(t, err)
	err = cfg.DBConfig.Set("test string")
	assert.Nil(t, err)
	lg, err := logger.NewLogger()
	assert.Nil(t, err)
	bd := mock_storage.NewMockdataBaser(ctrl)
	bd.EXPECT().Ping().Return(nil)
	bd.EXPECT().Get(gomock.Any()).Return("12345", nil).Times(1)
	bd.EXPECT().GetAllURLs(gomock.Any()).Return([]models.URLWithShort{{
		Shortcode: "1",
		URL:       "2",
	}}, nil).Times(1)
	bd.EXPECT().Get(gomock.Any()).Return("", fmt.Errorf("test error")).Times(1)
	st, err := NewStorage(cfg, lg, bd)
	assert.Nil(t, err)
	assert.NotNil(t, st)
	assert.True(t, st.Ping())
	url, err := st.dbGetURL("55555")
	assert.Nil(t, err)
	assert.Equal(t, "12345", url)
	_, err = st.dbGetURL("55555")
	assert.NotNil(t, err)
	ls, err := st.GetUserURLs("userID")
	assert.Nil(t, err)
	assert.True(t, len(ls) > 0)
}
