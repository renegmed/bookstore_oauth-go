package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	//"github.com/mercadolibre/golang-restclient/rest"
	resty "github.com/go-resty/resty/v2"
	"github.com/renegmed/bookstore_utils-go/rest_errors"
)

const (
	headerXPublic   = "X-Public"
	headerXClientId = "X-Client-Id"
	headerXCallerId = "X-Caller-Id"

	paramAccessToken = "access_token"
)

// var (
// 	// oauthRestClient = rest.RequestBuilder{
// 	// 	BaseURL: "http://localhost:8080",
// 	// 	Timeout: 200 * time.Millisecond,
// 	// }

// 	oauthRestClient = resty.New()
// )

type accessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	ClientId int64  `json:"client_id"`
}

type AccessTokenService struct {
	Client *resty.Client
}

func (ats *AccessTokenService) IsPublic() bool {
	if ats.Client.R() == nil {
		return true
	}

	//log.Println("+++++ header:", ats.Client.Header.Get(headerXPublic))

	return ats.Client.Header.Get(headerXPublic) == "true"
}

func (ats *AccessTokenService) GetCallerId() int64 {
	if ats.Client.R() == nil {
		return 0
	}
	callerId, err := strconv.ParseInt(ats.Client.R().Header.Get(headerXCallerId), 10, 64)
	if err != nil {
		return 0
	}
	return callerId
}

func (ats *AccessTokenService) GetClientId() int64 {
	if ats.Client.R() == nil {
		return 0
	}
	clientId, err := strconv.ParseInt(ats.Client.R().Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}
	return clientId
}

func (ats *AccessTokenService) AuthenticateRequest() rest_errors.RestErr {
	if ats.Client.R() == nil {
		return nil
	}

	ats.CleanRequest()

	//accessTokenId := strings.TrimSpace(ats.Client.R().URL.Query().Get(paramAccessToken))
	accessTokenId := strings.TrimSpace(ats.Client.R().QueryParam.Get(paramAccessToken))
	if accessTokenId == "" {
		return nil
	}

	at, err := ats.GetAccessToken(accessTokenId)
	if err != nil {
		if err.Status() == http.StatusNotFound {
			return nil
		}
		return err
	}
	ats.Client.R().Header.Add(headerXClientId, fmt.Sprintf("%v", at.ClientId))
	ats.Client.R().Header.Add(headerXCallerId, fmt.Sprintf("%v", at.UserId))
	return nil
}

func (ats *AccessTokenService) CleanRequest() {
	if ats.Client.R() == nil {
		return
	}
	ats.Client.R().Header.Del(headerXClientId)
	ats.Client.R().Header.Del(headerXCallerId)
}

func (ats *AccessTokenService) GetAccessToken(accessTokenId string) (*accessToken, rest_errors.RestErr) {

	response, err := ats.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "appication/json").
		Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenId))

	//fmt.Println(".....response:", response)

	if err != nil {
		return nil, rest_errors.NewInternalServerError("error on trying to get access token", err)
	}
	if response == nil || response.RawResponse == nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to get access token",
			errors.New("network timeout"))
	}

	if response.RawResponse.StatusCode > 299 {
		restErr, err := rest_errors.NewRestErrorFromBytes(response.Body())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to get access token", err)
		}
		return nil, restErr
	}

	var at accessToken
	if err := json.Unmarshal(response.Body(), &at); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal access token response",
			errors.New("error processing json"))
	}
	return &at, nil
}
