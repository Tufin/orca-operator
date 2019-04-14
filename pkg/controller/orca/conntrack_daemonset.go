package orca

import (
	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getConntrackDaemonset(cr *appv1alpha1.Orca) *appsv1.DaemonSet {
	labels := map[string]string{
		"app": "conntrack",
	}

	var selector = metav1.LabelSelector{
		MatchLabels: labels,
	}

	const name = "conntrack"
	const image = "docker.io/weaveworks/scope:1.10.1"

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector:        &selector,
			MinReadySeconds: 5,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name,
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: cr.Spec.Images["conntrack"],
							Args: []string{
								"'--mode=probe'",
								"'--probe-only'",
								"'--probe.kubernetes=true'",
								"'--probe.docker.bridge=docker0'",
								"'--probe.docker=true'",
								"'kite." + cr.Namespace + ":80'",
							},
							Env: []corev1.EnvVar{
								{
									Name:  "DOMAIN",
									Value: cr.Spec.Domain,
								},
							},
						},
					},
				},
			},
		},
	}
}
