package web

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path"

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
	mux.HandleFunc("/global/{part}", func(w http.ResponseWriter, r *http.Request) {
		base := path.Base(r.URL.Path)
		context := r.URL.Query().Get("context")
		http.Redirect(w, r, path.Join("/", context, base), http.StatusMovedPermanently)
	})
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
