package main

import (
	"context"
	"github.com/acorn-io/baaah"
	"github.com/acorn-io/baaah/pkg/log"
	"github.com/acorn-io/baaah/pkg/router"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"time"
)

const (
	reaperAnnotation = "ephemeral-reaper"
)

func main() {
	r, err := baaah.DefaultRouter("ephem-reaper", scheme.Scheme)
	if err != nil {
		panic(err)
	}

	r.Backend()

	r.Type(new(corev1.Pod)).HandlerFunc(handlePod)

	ctx := context.Background()
	err = r.Start(ctx)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	// Background context was canceled. Do whatever cleanup necessary.
}

func handlePod(req router.Request, resp router.Response) error {
	// Act on the pod, which can be retrieved via req.Object
	// If you need to create new objects based on the deployment, use resp.Objects()
	// A req.Client and req.Ctx are provided for instances when you need to directly interact with the Kubernetes API.
	pod := req.Object.(*corev1.Pod)

	if len(pod.Spec.EphemeralContainers) > 0 {
		log.Infof("Ephemeral containers found in pod %s", pod.Name)
		// check for annotation
		if _, ok := pod.Annotations[reaperAnnotation]; !ok {
			log.Infof("Marking pod %s for deletion in 1 hour", pod.Name)

			// Annotate the pod to indicate it is marked for deletion 1 hour from now
			timeToReap := time.Now().Add(time.Hour).Format(time.RFC3339)
			pod.SetAnnotations(map[string]string{reaperAnnotation: timeToReap})
			pod.SetDeletionTimestamp(&metav1.Time{Time: time.Now().Add(time.Hour)})
			return req.Client.Update(req.Ctx, pod)
		}
	}

	return nil
}
