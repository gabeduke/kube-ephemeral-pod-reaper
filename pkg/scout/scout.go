package scout

import (
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/labels"
	"time"
)

// NewScoutCmd creates and returns a new Scout command
func NewScoutCmd() *cobra.Command {
	var annotations string
	var duration time.Duration
	var name string
	var selector string

	var scoutCmd = &cobra.Command{
		Use:   "scout",
		Short: "Scout is a controller for marking ephemeral containers in Kubernetes for reaping",
		Long: `Scout watches Kubernetes pods for ephemeral containers and annotates them with 
               an expiration time. This is part of the kube-ephemeral-pod-reaper project, 
               where this controller marks containers for deletion.`,
		Run: func(cmd *cobra.Command, args []string) {
			var labelSelector labels.Selector
			var err error

			// Parse the label selector string into a labels.Selector
			if selector == "" {
				labelSelector = labels.Everything()
			} else {
				labelSelector, err = labels.Parse(selector)
				cobra.CheckErr(err)
			}

			cfg := Config{
				Annotation: annotations,
				Duration:   duration,
				Name:       name,
				Selector:   labelSelector,
			}

			controller := Controller{cfg: cfg}
			controller.Run()
		},
	}

	scoutCmd.Flags().StringVarP(&annotations, "annotations", "a", "ephemeral.reaper.leetserve.com/expiration-time", "Expiry Annotations to add to pods with ephemeral containers")
	scoutCmd.Flags().DurationVarP(&duration, "duration", "d", time.Hour, "Duration to wait before reaping")
	scoutCmd.Flags().StringVarP(&selector, "selector", "s", "", "Selector for pods to watch")
	scoutCmd.Flags().StringVarP(&name, "name", "n", "reaper-scout", "Name of the controller")

	return scoutCmd
}
