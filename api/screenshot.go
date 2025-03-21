package api

import (
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego"
	"github.com/go-fuego/fuego/option"
	"github.com/go-fuego/fuego/param"
	"github.com/yyewolf/rodent/mischief"
)

type ScreenshotRepository struct {
	mischief *mischief.Mischief
	logger   *slog.Logger
}

func NewScreenshotRepository(mischief *mischief.Mischief, logger *slog.Logger) *ScreenshotRepository {
	return &ScreenshotRepository{
		mischief: mischief,
		logger:   logger,
	}
}

func (s *ScreenshotRepository) Group() string {
	return "/screenshot"
}

var optionReturnsPNG = func(br *fuego.BaseRoute) {
	response := openapi3.NewResponse()
	response.WithDescription("Generated image")
	response.WithContent(openapi3.NewContentWithSchema(nil, []string{"image/png"}))
	br.Operation.AddResponse(200, response)
}

func (s *ScreenshotRepository) Register(server *fuego.Server) {
	fuego.GetStd(server, "", s.takeScreenshot,
		optionReturnsPNG,
		option.Description("Take a screenshot of the provided url."),
		option.Query("url", "The website to take a screenshot of", param.Example("example", "https://google.com")),
	)
}

func (s *ScreenshotRepository) takeScreenshot(writer http.ResponseWriter, req *http.Request) {
	unsafeUrl := req.URL.Query().Get("url")

	parsedUrl, err := url.Parse(unsafeUrl)
	if err != nil {
		s.logger.Error("error while parsing URL", slog.Any("error", err))
		http.Error(writer, "invalid URL", http.StatusBadRequest)
		return
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		err := errors.New("invalid URL scheme")
		s.logger.Error("error while parsing URL", slog.Any("error", err))
		http.Error(writer, "invalid URL scheme", http.StatusBadRequest)
		return
	}

	if parsedUrl.Port() != "" {
		err := errors.New("URL should not contain a port")
		s.logger.Error("error while parsing URL", slog.Any("error", err))
		http.Error(writer, "URL should not contain a port", http.StatusBadRequest)
		return
	}

	bytes, err := s.mischief.TakeScreenshot(parsedUrl.String())
	if err != nil {
		if errors.Is(err, mischief.ErrGettingBrowser) {
			http.Error(writer, "error while getting browser", http.StatusRequestTimeout)
			return
		}

		s.logger.Error("error while taking screenshot", slog.Any("error", err))
		http.Error(writer, "error while taking screenshot", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "image/png")
	_, err = writer.Write(bytes)
	if err != nil {
		s.logger.Error("error while writing response", slog.Any("error", err))
		http.Error(writer, "error while writing response", http.StatusInternalServerError)
		return
	}
}

var _ Repository = &ScreenshotRepository{}
