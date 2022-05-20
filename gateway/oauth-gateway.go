package gateway

// mockgen -source=gateway/oauth-gateway.go -destination=mock_gateway/oauth-gateway.go

import (
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type OauthGateway interface {
	SearchProvider(*gin.Context) (*oidc.Provider, error)
	RequestTokenEndpoint(oauth2Config oauth2.Config, ctx *gin.Context, code string) (*oauth2.Token, error)
	VerifyIDToken(ctx *gin.Context, provider *oidc.Provider, rawIDToken string) (*oidc.IDToken, error)
}

type oauthGateway struct{}

func NewOauthGateway() OauthGateway {
	return &oauthGateway{}
}

// googleの認可エンドポイントのurlやtokenエンドポイントのurlやid_tokenの署名の公開鍵のurl等を取りに行ってくれる
func (g *oauthGateway) SearchProvider(ctx *gin.Context) (*oidc.Provider, error) {
	return oidc.NewProvider(ctx.Request.Context(), os.Getenv("GOOGLE_OAUTH_URL"))
}

func (g *oauthGateway) RequestTokenEndpoint(oauth2Config oauth2.Config, ctx *gin.Context, code string) (*oauth2.Token, error) {
	return oauth2Config.Exchange(ctx.Request.Context(), code)
}

// id_tokenの検証 jwtの署名の公開鍵を取りに行く為gatewayに置く
func (g *oauthGateway) VerifyIDToken(ctx *gin.Context, provider *oidc.Provider, rawIDToken string) (*oidc.IDToken, error) {
	verifier := provider.Verifier(&oidc.Config{ClientID: os.Getenv("CLIENT_ID")})
	return verifier.Verify(ctx, rawIDToken)
}
