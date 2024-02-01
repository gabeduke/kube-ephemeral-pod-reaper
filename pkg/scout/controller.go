package scout

import (
	"context"
	"fmt"
	"github.com/acorn-io/baaah"
	"github.com/acorn-io/baaah/pkg/log"
	"github.com/acorn-io/baaah/pkg/router"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes/scheme"
	"time"
)

type Config struct {
	Name       string
	Annotation string
	Duration   time.Duration
	Selector   labels.Selector
}

type Controller struct {
	cfg Config
}

func (c *Controller) Run() {
	log.Infof = logrus.Infof
	log.Debugf = logrus.Debugf

	r, err := baaah.DefaultRouter(c.cfg.Name, scheme.Scheme)
	if err != nil {
		panic(err)
	}

	if c.cfg.Selector != nil {
		r.Selector(c.cfg.Selector)
	}

	r.Type(new(corev1.Pod)).HandlerFunc(c.handlePod)

	ctx := context.Background()

	err = r.Start(ctx)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	// Background context was canceled. Do whatever cleanup necessary.
}

func (c *Controller) handlePod(req router.Request, resp router.Response) error {
	// Act on the pod, which can be retrieved via req.Object
	// If you need to create new objects based on the deployment, use resp.Objects()
	// A req.Client and req.Ctx are provided for instances when you need to directly interact with the Kubernetes API.
	log.Debugf("Handling pod %s", req.Object.GetName())
	pod := req.Object.(*corev1.Pod)

	if len(pod.Spec.EphemeralContainers) > 0 {
		log.Debugf("Ephemeral containers found in pod %s", pod.Name)
		// check for annotation
		if _, ok := pod.Annotations[c.cfg.Annotation]; !ok {
			log.Infof("Marking pod %s for deletion in 1 hour", pod.Name)

			// Annotate the pod to indicate it is marked for deletion 1 hour from now
			timeToReap := time.Now().Add(c.cfg.Duration).Format(time.RFC3339)
			pod.SetAnnotations(map[string]string{c.cfg.Annotation: timeToReap})

			err := req.Client.Update(req.Ctx, pod)
			if err != nil {
				return fmt.Errorf("failed to update pod %s: %w", pod.Name, err)
			}
			return nil
		}
	}

	return nil
}
