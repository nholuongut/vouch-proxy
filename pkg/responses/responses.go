/*

Copyright 2020 The Vouch Proxy Authors.
Use of this source code is governed by The MIT License (MIT) that
can be found in the LICENSE file. Software distributed under The
MIT License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
OR CONDITIONS OF ANY KIND, either express or implied.

*/

package responses

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/nholuongut/vouch-proxy/pkg/cfg"
	"github.com/nholuongut/vouch-proxy/pkg/cookie"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

// Index variables passed to index.tmpl
type Index struct {
	Msg      string
	TestURLs []string
	Testing  bool
}

var (
	indexTemplate *template.Template
	errorTemplate *template.Template
	log           *zap.SugaredLogger
	fastlog       *zap.Logger

	errNotAuthorized = errors.New("not authorized")
)

// Configure see main.go configure()
func Configure() {
	log = cfg.Logging.Logger
	fastlog = cfg.Logging.FastLogger

	log.Debugf("responses.Configure() attempting to parse templates with cfg.RootDir: %s", cfg.RootDir)
	indexTemplate = template.Must(template.ParseFiles(filepath.Join(cfg.RootDir, "templates/index.tmpl")))

}

// RenderIndex render the response as an HTML page, mostly used in testing
func RenderIndex(w http.ResponseWriter, msg string) {
	if err := indexTemplate.Execute(w, &Index{Msg: msg, TestURLs: cfg.Cfg.TestURLs, Testing: cfg.Cfg.Testing}); err != nil {
		log.Error(err)
	}
}

// renderError html error page
// something terse for the end user
func renderError(w http.ResponseWriter, msg string) {
	log.Debugf("rendering error for user: %s", msg)
	if err := indexTemplate.Execute(w, &Index{Msg: msg}); err != nil {
		log.Error(err)
	}
}

// OK200 returns "200 OK"
func OK200(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("200 OK\n"))
	if err != nil {
		log.Error(err)
	}
}

// Redirect302 redirect to the specificed rURL
func Redirect302(w http.ResponseWriter, r *http.Request, rURL string) {
	if cfg.Cfg.Testing {
		cfg.Cfg.TestURLs = append(cfg.Cfg.TestURLs, rURL)
		RenderIndex(w, "302 redirect to: "+rURL)
		return
	}
	http.Redirect(w, r, rURL, http.StatusFound)
}

// Error400 Bad Request
func Error400(w http.ResponseWriter, r *http.Request, e error) {
	log.Error(e)
	cookie.ClearCookie(w, r)
	w.Header().Set(cfg.Cfg.Headers.Error, e.Error())
	w.WriteHeader(http.StatusBadRequest)
	addErrandCancelRequest(r)
	renderError(w, "400 Bad Request")
}

// Error401 Unauthorized the standard error
// this is captured by nginx, which converts the 401 into 302 to the login page
func Error401(w http.ResponseWriter, r *http.Request, e error) {
	log.Error(e)
	addErrandCancelRequest(r)
	cookie.ClearCookie(w, r)
	w.Header().Set(cfg.Cfg.Headers.Error, e.Error())
	http.Error(w, e.Error(), http.StatusUnauthorized)
	// renderError(w, "401 Unauthorized")
}

// Error403 Forbidden
// if there's an error during /auth or if they don't pass validation in /auth
func Error403(w http.ResponseWriter, r *http.Request, e error) {
	log.Error(e)
	addErrandCancelRequest(r)
	cookie.ClearCookie(w, r)
	w.Header().Set(cfg.Cfg.Headers.Error, e.Error())
	w.WriteHeader(http.StatusForbidden)
	renderError(w, "403 Forbidden")
}

// cfg.ErrCtx is tested by `jwtmanager.JWTCacheHandler`
func addErrandCancelRequest(r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	ctx = context.WithValue(ctx, cfg.ErrCtxKey, true)
	*r = *r.Clone(ctx)
	cancel() // we're done
	return
}
