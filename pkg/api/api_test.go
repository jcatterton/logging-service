package api

import (
	"errors"
	"logging-service/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"logging-service/pkg/testhelper/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestApi_CheckHealth_ShouldReturn500IfUnableToConnecteToDatabase(t *testing.T) {
	handler := &mocks.Handler{}
	handler.On("Ping", mock.Anything).Return(errors.New("test"))

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(checkHealth(handler))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, 500, recorder.Code)
}

func TestApi_CheckHealth_ShouldReturn200OnSuccess(t *testing.T) {
	handler := &mocks.Handler{}
	handler.On("Ping", mock.Anything).Return(nil)

	req, err := http.NewRequest(http.MethodGet, "/health", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(checkHealth(handler))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, 200, recorder.Code)
}

func TestApi_GetLogs_ShouldReturn500OnDaoError(t *testing.T) {
	handler := &mocks.Handler{}
	handler.On("GetLogs", mock.Anything).Return(nil, errors.New("test"))

	req, err := http.NewRequest(http.MethodGet, "/logs", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getLogs(handler))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, 500, recorder.Code)
}

func TestApi_GetLogs_ShouldReturn200OnSuccess(t *testing.T) {
	handler := &mocks.Handler{}
	handler.On("GetLogs", mock.Anything).Return([]models.Log{}, nil)

	req, err := http.NewRequest(http.MethodGet, "/logs", nil)
	require.Nil(t, err)

	recorder := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(getLogs(handler))
	httpHandler.ServeHTTP(recorder, req)
	require.Equal(t, 200, recorder.Code)
}
