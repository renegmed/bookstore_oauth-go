package oauth

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"

	//"github.com/mercadolibre/golang-restclient/rest"

	"github.com/jarcoal/httpmock"
)

func TestMain(m *testing.M) {
	fmt.Println("about to start oauth tests")

	//resty.StartMockupServer()

	os.Exit(m.Run())
}

func TestOauthConstants(t *testing.T) {
	assert.EqualValues(t, "X-Public", headerXPublic)
	assert.EqualValues(t, "X-Client-Id", headerXClientId)
	assert.EqualValues(t, "X-Caller-Id", headerXCallerId)
	assert.EqualValues(t, "access_token", paramAccessToken)
}

func TestIsPublicNilRequest(t *testing.T) {
	//defer httpmock.DeactivateAndReset()
	client := resty.New()
	ats := &AccessTokenService{client}
	ats.Client.SetHeader("X-Public", "true")

	// log.Println("... isPublic", ats.IsPublic())

	assert.True(t, ats.IsPublic())
}

func TestIsPublicNoError(t *testing.T) {
	rst := resty.New()
	ats := &AccessTokenService{rst}

	header := http.Header{}
	header.Set("X-Public", "true")
	// request := http.Request{
	// 	Header: header,
	// }

	assert.False(t, ats.IsPublic())
	//request.Header.Add("X-Public", "true")

	ats.Client.Header = header
	assert.True(t, ats.IsPublic())
}

func TestGetCallerIdNilRequest(t *testing.T) {
	//TODO: Complete!
}

func TestGetCallerInvalidCallerFormat(t *testing.T) {
	//TODO: Complete!
}

func TestGetCallerNoError(t *testing.T) {
	//TODO: Complete!
}

func newResponder(s int, c string, ct string) httpmock.Responder {
	resp := httpmock.NewStringResponse(s, c)
	resp.Header.Set("Content-Type", ct)

	return httpmock.ResponderFromResponse(resp)
}

func TestGetAccessTokenValidRestclientResponse(t *testing.T) {

	defer httpmock.DeactivateAndReset()

	rst := resty.New()
	ats := &AccessTokenService{rst}

	httpmock.ActivateNonDefault(rst.GetClient())
	httpmock.RegisterResponder(
		"GET", "/oauth/access_token/AbC123",
		newResponder(200, `
			{"id": "12345", "user_id": 12345, "client_id": 56789}
		`, "application/json"))

	accessToken, err := ats.GetAccessToken("AbC123")
	if err != nil {
		t.Fail()
	}

	assert.NotNil(t, accessToken)
	assert.Equal(t, "12345", accessToken.Id)
	assert.Equal(t, int64(12345), accessToken.UserId)
	assert.Equal(t, int64(56789), accessToken.ClientId)
}

func TestGetAccessTokenInvalidUrlRestclientResponse(t *testing.T) {

	defer httpmock.DeactivateAndReset()

	rst := resty.New()
	ats := &AccessTokenService{rst}

	httpmock.ActivateNonDefault(rst.GetClient())
	httpmock.RegisterResponder(
		"GET", "/oauth/access_token/AbC123",
		newResponder(500, `{}`, "application/json"))

	accessToken, err := ats.GetAccessToken("AbC12345678")
	//log.Println("++++ access token", accessToken, " error status code", http.StatusInternalServerError)

	assert.Nil(t, accessToken)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error on trying to get access token", err.Message())
}

func TestGetAccessTokenInvalidResponseRestclientResponse(t *testing.T) {

	defer httpmock.DeactivateAndReset()

	rst := resty.New()
	ats := &AccessTokenService{rst}

	httpmock.ActivateNonDefault(rst.GetClient())
	httpmock.RegisterResponder(
		"GET", "/oauth/access_token/AbC123",
		newResponder(500, ``, "application/json"))

	accessToken, err := ats.GetAccessToken("AbC123")
	//log.Println("++++ access token", accessToken, " error status code", http.StatusInternalServerError)

	assert.Nil(t, accessToken)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "invalid error interface when trying to get access token", err.Message())
}

func TestGetAccessTokenInvalidUnmarshalResponseRestclientResponse(t *testing.T) {

	defer httpmock.DeactivateAndReset()

	rst := resty.New()
	ats := &AccessTokenService{rst}

	httpmock.ActivateNonDefault(rst.GetClient())
	httpmock.RegisterResponder(
		"GET", "/oauth/access_token/AbC123",
		newResponder(200, `{"id": "12345", "user_id: 12345, "client_id": 56789}`, "application/json")) // note invalid json structure

	accessToken, err := ats.GetAccessToken("AbC123")
	// log.Println("++++ access token", accessToken, " error status code", http.StatusInternalServerError)

	assert.Nil(t, accessToken)
	assert.NotNil(t, err)

	assert.Equal(t, http.StatusInternalServerError, err.Status())
	assert.EqualValues(t, "error when trying to unmarshal access token response", err.Message())
}
