package orca

import (
	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func getKiteService(cr *appv1alpha1.Orca) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kite",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "syslog",
					Protocol: "UDP",
					Port:     6062,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1514,
						StrVal: "1514",
					},
				},
				{
					Name:     "grpc",
					Protocol: "TCP",
					Port:     81,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 6061,
						StrVal: "6061",
					},
				}, {
					Name:     "http",
					Protocol: "TCP",
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 6060,
						StrVal: "6060",
					},
				},
			},
			Selector: labels,
		},
	}
}
