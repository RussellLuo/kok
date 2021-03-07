// Code generated by kok; DO NOT EDIT.
// github.com/RussellLuo/kok

package helloworld

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/RussellLuo/kok/pkg/codec/httpcodec"
)

type HTTPClient struct {
	codecs     httpcodec.Codecs
	httpClient *http.Client
	scheme     string
	host       string
	pathPrefix string
}

func NewHTTPClient(codecs httpcodec.Codecs, httpClient *http.Client, baseURL string) (*HTTPClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	return &HTTPClient{
		codecs:     codecs,
		httpClient: httpClient,
		scheme:     u.Scheme,
		host:       u.Host,
		pathPrefix: strings.TrimSuffix(u.Path, "/"),
	}, nil
}

func (c *HTTPClient) SayHello(ctx context.Context, name string) (message string, err error) {
	codec := c.codecs.EncodeDecoder("SayHello")

	path := "/messages"
	u := &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.pathPrefix + path,
	}

	reqBody := struct {
		Name string `json:"name"`
	}{
		Name: name,
	}
	reqBodyReader, headers, err := codec.EncodeRequestBody(&reqBody)
	if err != nil {
		return "", err
	}

	_req, err := http.NewRequest("POST", u.String(), reqBodyReader)
	if err != nil {
		return "", err
	}

	for k, v := range headers {
		_req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(_req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode <= http.StatusNoContent {
		respBody := &SayHelloResponse{}
		err := codec.DecodeSuccessResponse(resp.Body, respBody.Body())
		if err != nil {
			return "", err
		}
		return respBody.Message, nil
	} else {
		var respErr error
		err := codec.DecodeFailureResponse(resp.Body, &respErr)
		if err == nil {
			err = respErr
		}
		return "", err
	}
}
