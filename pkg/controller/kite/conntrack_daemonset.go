package kite

import (
	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/app/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getConntrackDaemonset(cr *appv1alpha1.Kite) *appsv1.DaemonSet {
	labels := map[string]string{
		"app": cr.Name,
	}

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
	}
}
