package leaderelection

import (
	"context"
	"time"

	"github.com/gotway/gotway/pkg/kubernetes/controller"
	"github.com/gotway/gotway/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type Config struct {
	Identity           string
	LeaseLockName      string
	LeaseLockNamespace string
	LeaseDuration      time.Duration
	RenewDeadline      time.Duration
	RetryPeriod        time.Duration
	OnStartedLeading   func(context.Context)
	OnStoppedLeading   func()
}

type LeaderElection struct {
	ctrl      *controller.Controller
	clientset *kubernetes.Clientset
	logger    log.Logger
	config    Config
}

func (l *LeaderElection) Start(ctx context.Context) {
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      l.config.LeaseLockName,
			Namespace: l.config.LeaseLockNamespace,
		},
		Client: l.clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: l.config.Identity,
		},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   l.config.LeaseDuration,
		RenewDeadline:   l.config.RenewDeadline,
		RetryPeriod:     l.config.RetryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: l.config.OnStartedLeading,
			OnStoppedLeading: l.config.OnStoppedLeading,
		},
	})
}

func New(ctrl *controller.Controller, clientset *kubernetes.Clientset,
	logger log.Logger, config Config) *LeaderElection {

	return &LeaderElection{
		ctrl:      ctrl,
		clientset: clientset,
		logger:    logger,
		config:    config,
	}
}
