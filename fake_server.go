package apimock

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FakeEndpoint struct {
	HttpMethod   string
	Path         string
	ResponseJson string
	ResponseCode int
}

type ExecutedEndpointCall struct {
	HttpMethod        string
	Path              string
	RequestBodyInJson string
	RequestHeaders    map[string][]string
}

type FakeApi struct {
	ExecutedEndpointCalls []ExecutedEndpointCall
	HttpTestServer        *httptest.Server
}

func (fakeApi *FakeApi) AddExecutedEndpointCall(item ExecutedEndpointCall) []ExecutedEndpointCall {
	fakeApi.ExecutedEndpointCalls = append(fakeApi.ExecutedEndpointCalls, item)
	return fakeApi.ExecutedEndpointCalls
}

func StartFakeAPIServer(t *testing.T, fakeEndpoints []*FakeEndpoint) *FakeApi {

	fakeApi := &FakeApi{
		ExecutedEndpointCalls: []ExecutedEndpointCall{},
	}

	fakeApi.HttpTestServer = httptest.NewServer(http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		t.Logf("Method: %v", request.Method)
		t.Logf("Path: %v", request.URL.Path)
		for _, endpoint := range fakeEndpoints {

			if request.Method == endpoint.HttpMethod &&
				request.URL.Path == endpoint.Path {

				responseWriter.WriteHeader(endpoint.ResponseCode)
				_, _ = responseWriter.Write([]byte(endpoint.ResponseJson))

				body, _ := ioutil.ReadAll(request.Body)
				fakeApi.AddExecutedEndpointCall(ExecutedEndpointCall{
					HttpMethod:        endpoint.HttpMethod,
					Path:              endpoint.Path,
					RequestBodyInJson: string(body),
					RequestHeaders:    request.Header,
				})

				return
			}
		}
	}))

	return fakeApi
}

func (fakeApi *FakeApi) stop() {
	fakeApi.HttpTestServer.Close()
}
