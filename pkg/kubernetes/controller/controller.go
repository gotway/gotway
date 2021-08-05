package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	crdv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1"
	clientsetv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1/apis/clientset/versioned"
	informersv1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1/apis/informers/externalversions"
	"github.com/gotway/gotway/pkg/log"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Options struct {
	Namespace string
}

type Controller struct {
	options             Options
	ingresshttpInformer cache.SharedIndexInformer
	queue               workqueue.RateLimitingInterface
	logger              log.Logger
}

func (c *Controller) Run(ctx context.Context) error {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	c.logger.Info("starting controller")

	c.logger.Info("starting informer")
	go c.ingresshttpInformer.Run(ctx.Done())

	c.logger.Info("waiting for informer caches to sync")
	if !cache.WaitForCacheSync(ctx.Done(), c.ingresshttpInformer.HasSynced) {
		err := errors.New("failed to wait for informers caches to sync")
		utilruntime.HandleError(err)
		return err
	}
	c.logger.Info("controller ready")

	<-ctx.Done()
	c.logger.Info("stopping controller")

	return nil
}

func (c *Controller) List() ([]*crdv1alpha1.IngressHTTP, error) {
	var ingresses []*crdv1alpha1.IngressHTTP
	for _, obj := range c.ingresshttpInformer.GetIndexer().List() {
		if ingress, ok := obj.(*crdv1alpha1.IngressHTTP); ok {
			ingresses = append(ingresses, ingress)
			continue
		}
		return nil, fmt.Errorf("unexpected object %v", obj)
	}
	return ingresses, nil
}

func New(
	options Options,
	ingresshttpClientSet clientsetv1alpha1.Interface,
	logger log.Logger,
) *Controller {

	informerFactory := informersv1alpha1.NewSharedInformerFactory(ingresshttpClientSet, 10*time.Second)
	ingresshttpInformer := informerFactory.Gotway().V1alpha1().IngressHTTPs().Informer()

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	return &Controller{
		ingresshttpInformer: ingresshttpInformer,
		queue:               queue,
		options:             options,
		logger:              logger,
	}
}
