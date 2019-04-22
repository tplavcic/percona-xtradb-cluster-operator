package perconaxtradbcluster

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "github.com/percona/percona-xtradb-cluster-operator/pkg/apis/pxc/v1alpha1"
)

func (r *ReconcilePerconaXtraDBCluster) reconcilePods(namespace string, sfs api.StatefulApp) error {
	log.Info("PODS RECONCILING")

	list := corev1.PodList{}
	err := r.client.List(context.TODO(),
		&client.ListOptions{
			Namespace:     namespace,
			LabelSelector: labels.SelectorFromSet(sfs.Labels()),
		},
		&list,
	)
	if err != nil {
		return fmt.Errorf("get list: %v", err)
	}

	for _, pod := range list.Items {
		if pod.ObjectMeta.DeletionTimestamp != nil {
			// finalizers := []string{}
			for _, fnlz := range pod.GetFinalizers() {
				switch fnlz {
				case "pxc.percona.com/shutdown.pod.gracefully":
					log.Info("shutdown-pxc-pod-gracefully", "pod", pod.Name)
				}
			}
			pod.SetFinalizers(nil)
		}
	}

	return nil
}
