package matchingress

import (
	"net/http"
	"strings"

	httpError "github.com/gotway/gotway/internal/http/error"
	"github.com/gotway/gotway/internal/middleware"
	"github.com/gotway/gotway/internal/requestcontext"
	kubeCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	"github.com/gotway/gotway/pkg/log"

	crdv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1"
)

type matchIngress struct {
	kubeCtrl *kubeCtrl.Controller
	logger   log.Logger
}

func (m *matchIngress) MiddlewareFunc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug("match ingress")

		ingress, err := m.kubeCtrl.FindIngress(getIngressMatcher(r))
		if err != nil {
			httpError.Handle(err, w, m.logger)
			return
		}

		next.ServeHTTP(w, requestcontext.WithIngress(r, ingress))
	})
}

func getIngressMatcher(r *http.Request) kubeCtrl.IngressMatcher {
	return func(ingress *crdv1alpha1.IngressHTTP) bool {
		match := ingress.Spec.Match

		if match.Method != "" && match.Method != r.Method {
			return false
		}
		if match.Host != "" && match.Host != r.Host {
			return false
		}
		if match.Port != "" && match.Port != r.URL.Port() {
			return false
		}
		if match.Path != "" && match.Path != r.URL.RawPath {
			return false
		}
		if match.PathPrefix != "" && !strings.HasPrefix(r.URL.RawPath, match.PathPrefix) {
			return false
		}
		return true
	}
}

func New(
	kubeCtrl *kubeCtrl.Controller,
	logger log.Logger,
) middleware.Middleware {
	return &matchIngress{
		kubeCtrl: kubeCtrl,
		logger:   logger,
	}
}
