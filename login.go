package larki

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkauthen "github.com/larksuite/oapi-sdk-go/v3/service/authen/v1"
)

var jsApiReq = &larkcore.ApiReq{
	HttpMethod:                http.MethodGet,
	ApiPath:                   "https://open.feishu.cn/open-apis/jssdk/ticket/get",
	Body:                      nil,
	SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeTenant},
}

func (c *Client) GetWebAppUserAccessToken(ctx context.Context, code string) (*larkauthen.CreateOidcAccessTokenRespData, error) {
	resp, err := c.Authen.OidcAccessToken.Create(ctx,
		larkauthen.NewCreateOidcAccessTokenReqBuilder().
			Body(larkauthen.NewCreateOidcAccessTokenReqBodyBuilder().
				Code(code).GrantType("authorization_code").Build()).Build())
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) GetJsApiTicket(ctx context.Context) (*JsApiTicket, error) {
	resp, err := c.Do(ctx, jsApiReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %d", "Error fetching js api ticket", resp.StatusCode)
	}

	var ticketResp JsApiTicketResponse
	if err = sonic.Unmarshal(resp.RawBody, &ticketResp); err != nil {
		return nil, err
	}

	if ticketResp.Code != 0 {
		return nil, fmt.Errorf("%s: [%d] %s", "Error fetching js api ticket", ticketResp.Code, ticketResp.Message)
	}

	return &ticketResp.Data, nil
}

func (c *Client) GetJsAuthResponse(ticket, url, nocestr string) *JsApiAuthResponse {
	timestamp := time.Now().UnixMilli()

	return &JsApiAuthResponse{
		Appid: c.AppID,
		Signature: CestSign([]byte(fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
			ticket, nocestr, timestamp, url))),
		Noncestr:  nocestr,
		Timestamp: timestamp,
	}
}

func (c *Client) GetUserInfo(ctx context.Context, token string) (*larkauthen.GetUserInfoRespData, error) {
	resp, err := c.Authen.UserInfo.Get(ctx, larkcore.WithUserAccessToken(token))
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) GetMiniProgramUserAccessToken(ctx context.Context, code string) (*MiniProgramToken, error) {
	miniProgTokenReq := &larkcore.ApiReq{
		HttpMethod: http.MethodPost,
		ApiPath:    "https://open.feishu.cn/open-apis/mina/v2/tokenLoginValidate",
		Body: struct {
			Code string `json:"code"`
		}{
			Code: code,
		},
		SupportedAccessTokenTypes: []larkcore.AccessTokenType{larkcore.AccessTokenTypeApp},
	}

	resp, err := c.Do(ctx, miniProgTokenReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %d", "Error fetching mini program token", resp.StatusCode)
	}

	var miniProgTokenResp MiniProgramTokenResponse

	if err = sonic.Unmarshal(resp.RawBody, &miniProgTokenResp); err != nil {
		return nil, err
	}

	if miniProgTokenResp.Code != 0 {
		return nil, fmt.Errorf("%s: [%d] %s", "Error fetching mini program token", miniProgTokenResp.Code, miniProgTokenResp.Message)
	}

	return &miniProgTokenResp.Data, nil
}
