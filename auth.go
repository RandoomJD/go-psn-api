package psn

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Api struct {
	*Authenticator
	npsso  string
	tokens tokens
}

var authURL = must(url.Parse("https://ca.account.sony.com/api"))

type tokens struct {
	Access         string `json:"access_token"`
	Refresh        string `json:"refresh_token"`
	AccessExpires  int32  `json:"expires_in"`
	RefreshExpires int32  `json:"refresh_token_expires_in"`
}

var (
	ErrNPSSOEmpty  = errors.New("npsso is empty")
	ErrNPSSOLength = errors.New("npsso must be exactly 64 characters")
)

func validateNPSSO(npsso string) error {
	if npsso == "" {
		return ErrNPSSOEmpty
	}
	if len(npsso) != 64 {
		return ErrNPSSOLength
	}
	return nil
}

func (auth *Authenticator) Authenticate(ctx context.Context, npsso string) (*Api, error) {
	err := validateNPSSO(npsso)
	if err != nil {
		return nil, fmt.Errorf("invalid npsso: %v", err)
	}
	tokens, err := auth.authRequest(ctx, npsso)
	if err != nil {
		return nil, fmt.Errorf("can't do auth request: %w", err)
	}
	var out Api
	out.npsso = npsso
	out.tokens = tokens
	out.Authenticator = auth
	return &out, nil
}

func (api *Api) RefreashAccessToken(ctx context.Context) error {
	if api.tokens.Refresh == "" {
		return fmt.Errorf("refresh token is empty")
	}
	postValues := url.Values{}
	postValues.Add("scope", "psn:mobile.v1 psn:clientapp")
	postValues.Add("refresh_token", api.tokens.Refresh)
	postValues.Add("grant_type", "refresh_token")
	postValues.Add("token_format", "jwt")

	var postHeaders = authHeaders("", true)
	var tokens tokens
	url := authURL.JoinPath("authz/v3/oauth/token")
	err := api.post(ctx, postValues, url, postHeaders, &tokens)
	if err != nil {
		return fmt.Errorf("can't create new POST request %w: ", err)
	}
	api.tokens = tokens
	return nil
}

func (auth *Authenticator) authRequest(ctx context.Context, npsso string) (tokens, error) {
	getValues := url.Values{}
	getValues.Add("access_type", "offline")
	getValues.Add("app_context", "inapp_ios")
	getValues.Add("auth_ver", "v3")
	getValues.Add("cid", "60351282-8C5F-4D5E-9033-E48FEA973E11")
	getValues.Add("client_id", "ac8d161a-d966-4728-b0ea-ffec22f69edc")
	getValues.Add("darkmode", "true")
	getValues.Add("device_base_font_size", "10")
	getValues.Add("device_profile", "mobile")
	getValues.Add("duid", "0000000d0004008088347AA0C79542D3B656EBB51CE3EBE1")
	getValues.Add("elements_visibility", "no_aclink")
	getValues.Add("extraQueryParams", `{PlatformPrivacyWs1=minimal;}`)
	getValues.Add("no_captcha", "true")
	getValues.Add("redirect_uri", "com.playstation.PlayStationApp://redirect")
	getValues.Add("response_type", "code")
	getValues.Add("scope", "psn:mobile.v1 psn:clientapp")
	getValues.Add("service_entity", "urn:service-entity:psn")
	getValues.Add("service_logo", "ps")
	getValues.Add("smcid", "psapp:settings-entrance")
	getValues.Add("support_scheme", "sneiprls")
	getValues.Add("token_format", "jwt")
	getValues.Add("ui", "pr")

	uri := authURL.JoinPath("authz/v3/oauth/authorize")
	uri.RawQuery = getValues.Encode()

	nextUrl := uri.String()

	// not a best way to check redirect, refactor somewhere
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, nextUrl, nil)
	if err != nil {
		return tokens{}, fmt.Errorf("can't create new GET request: %w ", err)
	}

	req.Header = authHeaders(npsso, false)

	// create new httpclient with ability to check redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return tokens{}, fmt.Errorf("can't execute GET request: %w ", err)
	}

	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusFound {
		return tokens{}, fmt.Errorf("code doesn't retrieved from redirect")
	}

	nextUrl = resp.Header.Get("Location")
	parsed, err := url.ParseQuery(nextUrl)
	if err != nil {
		return tokens{}, fmt.Errorf("can't parse query: %w ", err)
	}

	if len(parsed["error_description"]) > 0 {
		return tokens{}, fmt.Errorf("can't authorize, error from psn: %s, check npsso (%s)", parsed["error_description"][0], npsso)
	}
	if len(parsed["com.playstation.PlayStationApp://redirect/?code"]) == 0 {
		return tokens{}, fmt.Errorf("can't get code")
	}

	code := parsed["com.playstation.PlayStationApp://redirect/?code"][0]

	postValues := url.Values{}
	postValues.Add("smcid", "psapp%3Asettings-entrance")
	postValues.Add("access_type", "offline")
	postValues.Add("code", code)
	postValues.Add("service_logo", "ps")
	postValues.Add("ui", "pr")
	postValues.Add("elements_visibility", "no_aclink")
	postValues.Add("redirect_uri", "com.playstation.PlayStationApp://redirect")
	postValues.Add("support_scheme", "sneiprls")
	postValues.Add("grant_type", "authorization_code")
	postValues.Add("darkmode", "true")
	postValues.Add("device_base_font_size", "10")
	postValues.Add("device_profile", "mobile")
	postValues.Add("app_context", "inapp_ios")
	postValues.Add("extraQueryParams", `{PlatformPrivacyWs1=minimal;}`)
	postValues.Add("token_format", "jwt")

	var postHeaders = authHeaders(npsso, true)

	var t tokens
	err = auth.post(ctx, postValues, authURL.JoinPath("authz/v3/oauth/token"), postHeaders, &t)
	if err != nil {
		return tokens{}, fmt.Errorf("can't create new POST request: %w ", err)
	}

	return t, nil
}

func authHeaders(npsso string, auth bool) http.Header {
	h := http.Header{}
	if npsso != "" {
		h.Add("Cookie", fmt.Sprintf("npsso=%s", npsso))
	}
	if auth {
		h.Add("Content-Type", "application/x-www-form-urlencoded")
		h.Add("Authorization", "Basic YWM4ZDE2MWEtZDk2Ni00NzI4LWIwZWEtZmZlYzIyZjY5ZWRjOkRFaXhFcVhYQ2RYZHdqMHY=")
	}
	return h
}
