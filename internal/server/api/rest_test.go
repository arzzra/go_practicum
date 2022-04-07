package api

import (
	"fmt"
	"github.com/arzzra/go_practicum/internal/server"
	"github.com/arzzra/go_practicum/internal/server/metric"
	"github.com/arzzra/go_practicum/internal/server/storage"
	localmock "github.com/arzzra/go_practicum/internal/server/storage/local/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func doTestRequest(t *testing.T, server *httptest.Server,
	method, path string) (int, string) {

	request, err := http.NewRequest(method, server.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	return response.StatusCode, string(responseBody)
}

func mockTestHandler(t *testing.T, mstorage storage.MetricStorage) *Handler {
	srv, err := server.MakeServer(mstorage)
	if err != nil {
		t.Fatal(err.Error())
	}

	h, err := MakeHandler(srv)
	if err != nil {
		t.Fatal(err.Error())
	}
	return h
}

func TestHandler_postMetric(t *testing.T) {
	type want struct {
		code int
	}
	tests := []struct {
		name    string
		uriPath string
		want    want
	}{
		// TODO: Add test cases.
		{
			name:    "positive test#1",
			uriPath: "/update/gauge/TEST1/0.12345",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:    "negative test#2",
			uriPath: "/update/notgauge/metric/44.4",
			want: want{
				code: http.StatusNotImplemented,
			},
		},
		{
			name:    "positive test#3",
			uriPath: "/update/counter/Alloc/42",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:    "positive test#4",
			uriPath: "/update/counter/Alloc/11",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:    "negative test#5",
			uriPath: "/update/counter/Alloc/44.4",
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockLocal := localmock.NewMockMetricStorage(mockCtrl)
	h := mockTestHandler(t, mockLocal)
	srv := httptest.NewServer(h.Router)
	defer srv.Close()

	mTest1, _ := metric.MakeMetricStruct("TEST1", "gauge", "0.12345")
	mTest3, _ := metric.MakeMetricStruct("Alloc", "counter", "42")
	mTest4, _ := metric.MakeMetricStruct("Alloc", "counter", "11")
	gomock.InOrder(
		mockLocal.EXPECT().UpdateMetric(mTest1).Return(nil),
		mockLocal.EXPECT().UpdateMetric(gomock.Any()).Times(0),
		mockLocal.EXPECT().UpdateMetric(mTest3).Return(nil),
		mockLocal.EXPECT().UpdateMetric(mTest4).Return(nil),
		mockLocal.EXPECT().UpdateMetric(gomock.Any()).Times(0),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, _ := doTestRequest(t, srv, http.MethodPost, tt.uriPath)
			assert.Equal(t, tt.want.code, statusCode)
		})
	}
}

func TestHandler_getMetric(t *testing.T) {
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name string
		path string
		want want
	}{
		{
			name: "positive test#1",
			path: "/value/gauge/TEST1",
			want: want{
				code: http.StatusOK,
				body: "0.12345",
			},
		},
		{
			name: "positive test#2",
			path: "/value/counter/Alloc",
			want: want{
				code: http.StatusOK,
				body: "53",
			},
		},
		{
			name: "negative test#3",
			path: "/value/asdf/metric2",
			want: want{
				code: http.StatusNotImplemented,
				body: "unknown type of metric: asdf\n",
			},
		},
		{
			name: "negative test#4",
			path: "/value/counter/nemo",
			want: want{
				code: http.StatusNotFound,
				body: "Metric nemo not found\n",
			},
		},
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockLocal := localmock.NewMockMetricStorage(mockCtrl)
	h := mockTestHandler(t, mockLocal)
	srv := httptest.NewServer(h.Router)
	defer srv.Close()

	mTest1, _ := metric.MakeMetricStruct("TEST1", "gauge", "0.12345")
	mTest2, _ := metric.MakeMetricStruct("Alloc", "counter", "53")

	gomock.InOrder(
		mockLocal.EXPECT().GetMetricFromStorage(
			gomock.Any(),
			gomock.Any(),
		).Return(&mTest1, nil),

		mockLocal.EXPECT().GetMetricFromStorage(
			gomock.Any(),
			gomock.Any(),
		).Return(&mTest2, nil),
		mockLocal.EXPECT().GetMetricFromStorage(
			gomock.Any(),
			gomock.Any(),
		).Times(0),
		mockLocal.EXPECT().GetMetricFromStorage(
			gomock.Any(),
			gomock.Any(),
		).Return(nil, fmt.Errorf("Metric nemo not found")),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusCode, body := doTestRequest(t, srv, http.MethodGet, tt.path)
			assert.Equal(t, tt.want.code, statusCode)
			assert.Equal(t, tt.want.body, body)
		})
	}
}
