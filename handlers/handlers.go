/*

Copyright 2020 The Vouch Proxy Authors.
Use of this source code is governed by The MIT License (MIT) that
can be found in the LICENSE file. Software distributed under The
MIT License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
OR CONDITIONS OF ANY KIND, either express or implied.

*/

package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"go.uber.org/zap"

	"github.com/nholuongut/vouch-proxy/pkg/cfg"
	"github.com/nholuongut/vouch-proxy/pkg/cookie"
	"github.com/nholuongut/vouch-proxy/pkg/providers/adfs"
	"github.com/nholuongut/vouch-proxy/pkg/providers/azure"
	"github.com/nholuongut/vouch-proxy/pkg/providers/common"
	"github.com/nholuongut/vouch-proxy/pkg/providers/github"
	"github.com/nholuongut/vouch-proxy/pkg/providers/google"
	"github.com/nholuongut/vouch-proxy/pkg/providers/homeassistant"
	"github.com/nholuongut/vouch-proxy/pkg/providers/indieauth"
	"github.com/nholuongut/vouch-proxy/pkg/providers/nextcloud"
	"github.com/nholuongut/vouch-proxy/pkg/providers/openid"
	"github.com/nholuongut/vouch-proxy/pkg/providers/openstax"
	"github.com/nholuongut/vouch-proxy/pkg/structs"
)

// Provider each Provider must support GetuserInfo
type Provider interface {
	Configure()
	GetUserInfo(r *http.Request, user *structs.User, customClaims *structs.CustomClaims, ptokens *structs.PTokens) error
}

const (
	base64Bytes = 32
)

var (
	sessstore *sessions.CookieStore
	log       *zap.SugaredLogger
	fastlog   *zap.Logger
	provider  Provider
)

// Configure see main.go configure()
func Configure() {
	log = cfg.Logging.Logger
	fastlog = cfg.Logging.FastLogger
	// http://www.gorillatoolkit.org/pkg/sessions
	sessstore = sessions.NewCookieStore([]byte(cfg.Cfg.Session.Key))
	sessstore.Options.HttpOnly = cfg.Cfg.Cookie.HTTPOnly
	sessstore.Options.Secure = cfg.Cfg.Cookie.Secure
	sessstore.Options.SameSite = cookie.SameSite()
	sessstore.Options.MaxAge = 300 // give the user five minutes to log in at the IdP

	provider = getProvider()
	provider.Configure()
	common.Configure()
}

func getProvider() Provider {
	switch cfg.GenOAuth.Provider {
	case cfg.Providers.IndieAuth:
		return indieauth.Provider{}
	case cfg.Providers.ADFS:
		return adfs.Provider{}
	case cfg.Providers.Azure:
		return azure.Provider{}
	case cfg.Providers.HomeAssistant:
		return homeassistant.Provider{}
	case cfg.Providers.OpenStax:
		return openstax.Provider{}
	case cfg.Providers.Google:
		return google.Provider{}
	case cfg.Providers.GitHub:
		return github.Provider{PrepareTokensAndClient: common.PrepareTokensAndClient}
	case cfg.Providers.Nextcloud:
		return nextcloud.Provider{}
	case cfg.Providers.OIDC:
		return openid.Provider{}
	default:
		// shouldn't ever reach this since cfg checks for a properly configure `oauth.provider`
		log.Fatal("oauth.provider appears to be misconfigured, please check your config")
		return nil
	}
}
