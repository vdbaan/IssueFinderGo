/*
 *
 *  Copyright © 2021. Steven van der Baan <steven.vanderbaan@nccgroup.com>
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package web

import (
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	//rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"issuefinder/infra/config"
	"issuefinder/interface/web/api"
	"issuefinder/interface/web/support"
)

type WebServer interface {
	SetDebug(debug bool)
	RunLive(live bool)
	Start() error
	GetListeningAddress() string
}

type WebMainHandler struct {
	Config  config.Handler
	Address string
	Debug   bool
	Live    bool
}

func NewWebServer(address string, config config.Handler) WebServer {
	h := new(WebMainHandler)
	h.Address = address
	h.Config = config
	return h
}

func (h *WebMainHandler) GetListeningAddress() string {
	return h.Address
}

func (h *WebMainHandler) SetDebug(debug bool) {
	h.Debug = debug
}
func (h *WebMainHandler) RunLive(live bool) {
	h.Live = live
}

func (h *WebMainHandler) Start() error {
	e := echo.New()
	if h.Debug {
		e.Logger.SetLevel(log.DEBUG)
		e.Debug = true
		e.Use(middleware.Logger())
	} else {
		e.HideBanner = true
	}

	e.Use(support.NoCache())

	apiHandler := api.NewAPIHandler(h.Config)

	e.POST("/api/file/upload", apiHandler.UploadFile)
	e.POST("/api/findings", apiHandler.GetFindings)
	//e.POST("/api/findings/filter", apiHandler.GetFilteredFindings)
	e.GET("/api/filters", apiHandler.GetFilters)
	e.POST("/api/filters", apiHandler.Addfilter)
	//e.POST("/api/filters/risk", apiHandler.ShowSettings)
	e.GET("/api/reset", apiHandler.Reset)

	assetHandler := http.FileServer(getFileSystem(h.Live))

	e.GET("/*", echo.WrapHandler(assetHandler))

	// create a custom listener to capture the port we're listening on
	ln, _ := net.Listen("tcp", h.Address)
	h.Address = ln.Addr().String()
	e.Listener = ln

	// don't need to provide an address as we already have a custom listener
	return e.Start("")
}

//go:embed www
var embededFiles embed.FS

func getFileSystem(isLive bool) http.FileSystem {
	if isLive {
		fmt.Println("We're LIVE!!")
		s, _ := resolveAbsolutePathFromCaller("www", 1)
		return http.FS(os.DirFS(s))
	}

	fsys, err := fs.Sub(embededFiles, "www")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

// taken from go.rice, finds the folder of the function based on the stacktrace
func resolveAbsolutePathFromCaller(name string, nStackFrames int) (string, error) {
	_, callingGoFile, _, ok := runtime.Caller(nStackFrames)
	if !ok {
		return "", errors.New("couldn't find caller on stack")
	}
	// resolve to proper path
	pkgDir := filepath.Dir(callingGoFile)
	// fix for go cover
	const coverPath = "_test/_obj_test"
	if !filepath.IsAbs(pkgDir) {
		if i := strings.Index(pkgDir, coverPath); i >= 0 {
			pkgDir = pkgDir[:i] + pkgDir[i+len(coverPath):]            // remove coverPath
			pkgDir = filepath.Join(os.Getenv("GOPATH"), "src", pkgDir) // make absolute
		}
	}
	return filepath.Join(pkgDir, name), nil
}
