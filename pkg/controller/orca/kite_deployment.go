package orca

import (
	"strings"

	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kiteDeployment = "kite"
)

func getKiteDeployment(cr *appv1alpha1.Orca) *appsv1.Deployment {
	labels := map[string]string{
		"app": cr.Name,
	}

	var replicas int32 = 1

	var selector = metav1.LabelSelector{
		MatchLabels: labels,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kiteDeployment,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name,
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kite",
							Image: cr.Spec.Images["kite"],
							Env: []corev1.EnvVar{
								{
									Name:  "DOMAIN",
									Value: cr.Spec.Domain,
								},
								{
									Name:  "PROJECT",
									Value: cr.Spec.Project,
								},
								{
									Name:  "IGNORED_CONFIGMAPS",
									Value: strings.Join(cr.Spec.IngnoredConfigMaps, ","),
								},
								// Endpoints
								{
									Name:  "TUFIN_ORCA_URL",
									Value: cr.Spec.EndPoints["orca"],
								},
								{
									Name:  "TUFIN_GURU_URL",
									Value: cr.Spec.EndPoints["guru"],
								},
								{
									Name:  "TUFIN_DOCKER_REPO_URL",
									Value: cr.Spec.EndPoints["registry"],
								},
								// components
								{
									Name:  "TUFIN_INSTALL_DNS",
									Value: bts(cr.Spec.Components["dns"]),
								},
								{
									Name:  "TUFIN_INSTALL_CONNTRACK",
									Value: bts(cr.Spec.Components["conntrack"]),
								},
								{
									Name:  "TUFIN_INSTALL_SYSLOG",
									Value: bts(cr.Spec.Components["syslog"]),
								},
								{
									Name:  "TUFIN_INSTALL_ISTIO",
									Value: bts(cr.Spec.Components["istio"]),
								},
								{
									Name:  "TUFIN_INSTALL_DOCKER_PUSHER",
									Value: bts(cr.Spec.Components["docker_pusher"]),
								},
								{
									Name:  "TUFIN_INSTALL_KUBE_EVENTS_WATCHER",
									Value: bts(cr.Spec.Components["kube_events_watcher"]),
								},
								{
									Name:  "TUFIN_INSTALL_KUBE_EVENTS_WATCHER_NETWORK_POLICY",
									Value: bts(cr.Spec.Components["kube_events_watcher_network_policy"]),
								},
								// secrets
								{
									Name:      "CRT",
									ValueFrom: getSecretValue("orca-secrets", "guru-crt"),
								},
								{
									Name:      "TUFIN_DOCKER_REPO_USERNAME",
									ValueFrom: getSecretValue("orca-secrets", "docker-repo-username"),
								},
								{
									Name:      "TUFIN_DOCKER_REPO_PASSWORD",
									ValueFrom: getSecretValue("orca-secrets", "guru-api-key"),
								},
								{
									Name:      "API_KEY",
									ValueFrom: getSecretValue("orca-secrets", "guru-api-key"),
								},
							},
						},
					},
				},
			},
		},
	}
}

func bts(b bool) string {
	if b {
		return "TRUE"
	}
	return "FALSE"
}

func getSecretValue(name string, key string) *corev1.EnvVarSource {
	var optional = true

	return &corev1.EnvVarSource{
		SecretKeyRef: &corev1.SecretKeySelector{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: name,
			},
			Key:      key,
			Optional: &optional,
		},
	}
}
