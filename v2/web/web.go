package web

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/eduardooliveira/stLib/v2/web/comp"
)

type webHandler struct {
	l *slog.Logger
}

func New() (http.Handler, error) {
	wh := &webHandler{
		l: slog.With("module", "core-web"),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", R(wh.indexHandler))
	return mux, nil
}

func RenderIndex(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return comp.WrapperComponent(comp.WrapperModel{
		Main: templ.NopComponent,
	}).Render(context.Background(), f)

}
