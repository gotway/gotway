package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gotway/gotway/internal/cache"
	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/model"
	"github.com/gotway/gotway/internal/requestcontext"
	"github.com/gotway/gotway/pkg/log"

	kubeCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
)

type handler struct {
	kubeCtrl  *kubeCtrl.Controller
	cacheCtrl cache.Controller
	logger    log.Logger
}

func (h *handler) getIngresses(w http.ResponseWriter, r *http.Request) {
	ingresses, err := h.kubeCtrl.ListIngresses()
	if err != nil {
		httpError.Handle(err, w, h.logger)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ingresses)
}

func (h *handler) deleteCache(w http.ResponseWriter, r *http.Request) {
	decoded := json.NewDecoder(r.Body)

	var payload model.DeleteCache
	err := decoded.Decode(&payload)
	if err != nil {
		http.Error(w, model.ErrInvalidDeleteCache.Error(), http.StatusBadRequest)
		return
	}

	err = payload.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(payload.Paths) > 0 {
		err := h.cacheCtrl.DeleteCacheByPath(payload.Paths)
		if err != nil {
			if _, ok := err.(*model.ErrCachePathNotFound); ok {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if len(payload.Tags) > 0 {
		err := h.cacheCtrl.DeleteCacheByTags(payload.Tags)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) writeResponse(w http.ResponseWriter, r *http.Request) {
	res, err := requestcontext.GetResponse(r)
	if err != nil {
		httpError.Handle(err, w, h.logger)
		return
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		httpError.Handle(err, w, h.logger)
		return
	}

	h.logger.Debug("write response")
	for key, header := range res.Header {
		w.Header().Set(key, strings.Join(header[:], ","))
	}
	w.WriteHeader(res.StatusCode)
	w.Write(bytes)
}

func getServiceKey(r *http.Request) string {
	params := mux.Vars(r)
	return params["service"]
}

func newHandler(
	kubeCtrl *kubeCtrl.Controller,
	cacheController cache.Controller,
	logger log.Logger,
) *handler {

	return &handler{
		kubeCtrl:  kubeCtrl,
		cacheCtrl: cacheController,
		logger:    logger,
	}
}
