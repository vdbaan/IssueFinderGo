package web

import (
	"embed"
	_ "embed"
	"errors"
	"io/fs"
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
	Start() error
}

type WebMainHandler struct {
	Config  config.Handler
	Address string
	Debug   bool
}

func NewWebServer(address string, config config.Handler) WebServer {
	h := new(WebMainHandler)
	h.Address = address
	h.Config = config
	return h
}
func (h *WebMainHandler) SetDebug(debug bool) {
	h.Debug = debug
}

func (h *WebMainHandler) Start() error {
	e := echo.New()
	if h.Debug {
		e.Debug = true
	}
	e.Logger.SetLevel(log.DEBUG)
	e.Use(middleware.Logger())
	e.Use(support.NoCache())

	apiHandler := api.NewAPIHandler(h.Config)

	e.POST("/api/file/upload", apiHandler.UploadFile)
	e.POST("/api/findings", apiHandler.GetFindings)
	//e.POST("/api/findings/filter", apiHandler.GetFilteredFindings)
	e.GET("/api/filters", apiHandler.GetFilters)
	e.POST("/api/filters", apiHandler.Addfilter)
	//e.POST("/api/filters/risk", apiHandler.ShowSettings)
	e.GET("/api/reset", apiHandler.Reset)

	assetHandler := http.FileServer(getFileSystem())

	e.GET("/*", echo.WrapHandler(assetHandler))

	return e.Start(h.Address)
}

//go:embed www
var embededFiles embed.FS

// will be set during build
var runEmbedded = false

func getFileSystem() http.FileSystem {
	if !runEmbedded {
		s, _ := resolveAbsolutePathFromCaller("www", 1)
		return http.FS(os.DirFS(s))
	}

	fsys, err := fs.Sub(embededFiles, "")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

// taken from go.rice
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
