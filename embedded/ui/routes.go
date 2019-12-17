package ui

import (
	"html/template"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/tendermint/dex-demo/embedded/auth"
)

var tmpl *template.Template
var boxHdlr http.Handler
var uiPaths map[string]bool

func RegisterRoutes(_ context.CLIContext, r *mux.Router, _ *codec.Codec) {
	r.PathPrefix("/").HandlerFunc(uiHandler).Methods("GET")
}

func uiHandler(w http.ResponseWriter, r *http.Request) {
	if tmpl == nil {
		box := packr.NewBox("./public")
		boxHdlr = http.FileServer(box)
		tmplStr, err := box.FindString("/index.html")
		if err != nil {
			panic(err)
		}
		t, err := template.New("entry").Parse(tmplStr)
		if err != nil {
			panic(err)
		}
		tmpl = t
	}

	if _, ok := uiPaths[r.URL.Path]; ok {
		kb, err := auth.GetKBFromSession(r)
		var uexAddr string
		if err == nil {
			uexAddr = kb.GetAddr().String()
		}

		tmplState := TemplateState{
			UEXAddress: uexAddr,
		}
		err = tmpl.Execute(w, tmplState)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	boxHdlr.ServeHTTP(w, r)
}

type TemplateState struct {
	UEXAddress string
}

func init() {
	uiPaths = make(map[string]bool)
	uiPaths["/exchange"] = true
	uiPaths["/wallet"] = true
	uiPaths["/"] = true
}
