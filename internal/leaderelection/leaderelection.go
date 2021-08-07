package leaderelection

import (
	"context"
	"time"

	"github.com/gotway/gotway/internal/healthcheck"
	kubernetesCtrl "github.com/gotway/gotway/pkg/kubernetes/controller"
	"github.com/gotway/gotway/pkg/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Options struct {
	HealthCheckEnabled bool
	HAEnabled          bool
	Namespace          string
	NodeId             string
	LeaseLockName      string
	LeaseDuration      time.Duration
	RenewDeadline      time.Duration
	RetryPeriod        time.Duration
}

type Controller struct {
	options       Options
	healthCtrl    *healthcheck.Controller
	kubeCtrl      *kubernetesCtrl.Controller
	kubeClientSet *kubernetes.Clientset
	logger        log.Logger
}

func (c *Controller) Start(ctx context.Context) {
	if c.options.HAEnabled {
		c.logger.Info("starting HA controllers")
		c.startHA(ctx)
	} else {
		c.logger.Info("starting standalone controllers")
		c.startSingleNode(ctx)
	}
}

func (c *Controller) startHA(ctx context.Context) {
	if c.options == (Options{}) || !c.options.HAEnabled {
		c.logger.Fatal("HA config not set or not enabled")
	}

	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      c.options.LeaseLockName,
			Namespace: c.options.Namespace,
		},
		Client: c.kubeClientSet.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: c.options.NodeId,
		},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   c.options.LeaseDuration,
		RenewDeadline:   c.options.RenewDeadline,
		RetryPeriod:     c.options.RetryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				c.logger.Info("start leading")
				c.startSingleNode(ctx)
			},
			OnStoppedLeading: func() {
				c.logger.Info("stopped leading")
			},
			OnNewLeader: func(identity string) {
				if identity == c.options.NodeId {
					c.logger.Info("obtained leadership")
					return
				}
				c.logger.Infof("leader elected: '%s'", identity)
			},
		},
	})
}

func (c *Controller) startSingleNode(ctx context.Context) {
	if c.options.HealthCheckEnabled {
		go c.healthCtrl.Start(ctx)
	}
	go c.startKubernetesCtrl(ctx)
}

func (c *Controller) startKubernetesCtrl(ctx context.Context) {
	if err := c.kubeCtrl.Run(ctx); err != nil {
		c.logger.Fatal("error starting kubernetes controller", err)
	}
}

func NewController(
	options Options,
	healthCtrl *healthcheck.Controller,
	kubeCtrl *kubernetesCtrl.Controller,
	kubeClientSet *kubernetes.Clientset,
	logger log.Logger,
) *Controller {

	return &Controller{
		options:       options,
		healthCtrl:    healthCtrl,
		kubeCtrl:      kubeCtrl,
		kubeClientSet: kubeClientSet,
		logger:        logger,
	}
}
