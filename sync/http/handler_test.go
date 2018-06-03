package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mantzas/patron/encoding"
	"github.com/mantzas/patron/encoding/json"
	"github.com/mantzas/patron/sync"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_extractHeaders(t *testing.T) {
	assert := assert.New(t)
	r, err := http.NewRequest("GET", "/", nil)
	assert.NoError(err)
	r.Header.Set("KEY1", "VALUE1")
	r.Header.Set("KEY2", "VALUE2")
	m := extractHeaders(r)
	assert.Len(m, 2)
	assert.Equal("VALUE1", m["Key1"])
	assert.Equal("VALUE2", m["Key2"])
}

func Test_extractFields(t *testing.T) {
	assert := assert.New(t)
	r, err := http.NewRequest("GET", "/test?value1=1&value2=2", nil)
	assert.NoError(err)
	f := extractFields(r)
	assert.Len(f, 2)
	assert.Equal("1", f["value1"])
	assert.Equal("2", f["value2"])
}

func Test_determineEncoding(t *testing.T) {

	assert := assert.New(t)
	hdrContentJSON := make(map[string]string)
	hdrContentJSON[encoding.ContentTypeHeader] = json.ContentTypeCharset
	hdrEmptyHeader := make(map[string]string)
	hdrUnsupportedEncoding := make(map[string]string)
	hdrUnsupportedEncoding[encoding.ContentTypeHeader] = "application/xml"

	type args struct {
		hdr map[string]string
	}
	tests := []struct {
		name    string
		args    args
		decode  encoding.Decode
		encode  encoding.Encode
		wantErr bool
	}{
		{"content type json", args{hdrContentJSON}, json.Decode, json.Encode, false},
		{"empty header", args{hdrEmptyHeader}, nil, nil, true},
		{"unsupported encoding", args{hdrUnsupportedEncoding}, nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct, got, got1, err := determineEncoding(tt.args.hdr)
			if tt.wantErr {
				assert.Error(err)
				assert.Nil(got)
				assert.Nil(got1)
				assert.Empty(ct)
			} else {
				assert.NoError(err)
				assert.NotNil(got)
				assert.NotNil(got1)
				assert.Equal(json.ContentTypeCharset, ct)
			}
		})
	}
}

func Test_handleSuccess(t *testing.T) {
	assert := assert.New(t)
	get, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	post, err := http.NewRequest(http.MethodPost, "/", nil)
	assert.NoError(err)
	jsonRsp := sync.NewResponse(struct {
		Name    string
		Address string
	}{"Sotiris", "Athens"})
	jsonEncodeFailRsp := sync.NewResponse(struct {
		Name    chan bool
		Address string
	}{nil, "Athens"})

	type args struct {
		r   *http.Request
		rsp *sync.Response
		enc encoding.Encode
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
		wantErr        bool
	}{
		{"GET No Content success", args{get, nil, nil}, http.StatusNoContent, false},
		{"GET OK success", args{get, jsonRsp, json.Encode}, http.StatusOK, false},
		{"POST Created success", args{post, jsonRsp, json.Encode}, http.StatusCreated, false},
		{"Encode failure", args{post, jsonEncodeFailRsp, json.Encode}, http.StatusCreated, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rsp := httptest.NewRecorder()

			err := handleSuccess(rsp, tt.args.r, tt.args.rsp, tt.args.enc)
			if tt.wantErr {
				assert.Error(err)

			} else {
				assert.NoError(err)
				assert.Equal(tt.expectedStatus, rsp.Code)
			}
		})
	}
}

func Test_handleError(t *testing.T) {
	assert := assert.New(t)
	type args struct {
		err error
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{"bad request", args{&sync.ValidationError{}}, http.StatusBadRequest},
		{"unauthorized request", args{&sync.UnauthorizedError{}}, http.StatusUnauthorized},
		{"forbidden request", args{&sync.ForbiddenError{}}, http.StatusForbidden},
		{"not found error", args{&sync.NotFoundError{}}, http.StatusNotFound},
		{"service unavailable error", args{&sync.ServiceUnavailableError{}}, http.StatusServiceUnavailable},
		{"default error", args{errors.New("Test")}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsp := httptest.NewRecorder()
			handleError(rsp, tt.args.err)
			assert.Equal(tt.expectedCode, rsp.Code)
		})
	}
}

type testHandler struct {
	err  bool
	resp interface{}
}

func (th testHandler) Process(ctx context.Context, req *sync.Request) (*sync.Response, error) {
	if th.err {
		return nil, errors.New("TEST")
	}
	return sync.NewResponse(th.resp), nil
}

func Test_handler(t *testing.T) {
	assert := assert.New(t)

	errReq, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(err)
	req.Header.Set(encoding.ContentTypeHeader, json.ContentType)

	// success handling
	// failure handling
	type args struct {
		req *http.Request
		hnd sync.Processor
	}
	tests := []struct {
		name         string
		args         args
		expectedCode int
	}{
		{"unsupported content type", args{errReq, nil}, http.StatusUnsupportedMediaType},
		{"success handling", args{req, testHandler{false, "test"}}, http.StatusOK},
		{"error handling", args{req, testHandler{true, "test"}}, http.StatusInternalServerError},
		{"success handling failed due to encoding", args{req, testHandler{false, make(chan bool)}}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsp := httptest.NewRecorder()
			handler(tt.args.hnd).ServeHTTP(rsp, tt.args.req)
			assert.Equal(tt.expectedCode, rsp.Code)
		})
	}
}

func Test_prepareResponse(t *testing.T) {
	assert := assert.New(t)
	rsp := httptest.NewRecorder()
	prepareResponse(rsp, json.ContentTypeCharset)
	assert.Equal(json.ContentTypeCharset, rsp.Header().Get(encoding.ContentTypeHeader))
}
