package orca

import (
	"strings"

	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func getKiteDeployment(cr *appv1alpha1.Orca) *appsv1.Deployment {

	var replicas int32 = 1
	labels := GetLabels(app + "=" + kite)

	if kiteImg := os.Getenv("RELATED_IMAGES_TUFIN_KITE"); kiteImg != "" {
		cr.Spec.Images[kite] = kiteImg
	}

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kite,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: GetLabelSelector(labels),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      kite,
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: kite,
					Volumes: []corev1.Volume{
						{
							Name: kiteTempVolumeName,
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:            kite,
							Image:           cr.Spec.Images[kite],
							ImagePullPolicy: corev1.PullAlways,
							VolumeMounts: []corev1.VolumeMount{
								{Name: kiteTempVolumeName, MountPath: kiteTempVolumePath},
							},
							SecurityContext: &corev1.SecurityContext{
								Privileged: GetBoolRef(true),
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 6060}, {ContainerPort: 6061}, {ContainerPort: 6062, Protocol: corev1.ProtocolUDP},
							},
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
									Value: BoolToString(cr.Spec.Components["dns"]),
								},
								{
									Name:  "TUFIN_INSTALL_CONNTRACK",
									Value: BoolToString(cr.Spec.Components["conntrack"]),
								},
								{
									Name:  "TUFIN_INSTALL_SYSLOG",
									Value: BoolToString(cr.Spec.Components["syslog"]),
								},
								{
									Name:  "TUFIN_INSTALL_ISTIO",
									Value: BoolToString(cr.Spec.Components["istio"]),
								},
								{
									Name:  "TUFIN_INSTALL_DOCKER_PUSHER",
									Value: BoolToString(cr.Spec.Components["pusher"]),
								},
								{
									Name:  "TUFIN_INSTALL_KUBE_EVENTS_WATCHER",
									Value: BoolToString(cr.Spec.Components["watcher"]),
								},
								{
									Name:  "TUFIN_INSTALL_NETWORK_POLICY",
									Value: BoolToString(cr.Spec.Components["kube-network-policy"]),
								},
								// secrets
								{
									Name:      "CRT",
									ValueFrom: getSecretValue(kiteSecretsName, "guru-crt"),
								},
								{
									Name:      "TUFIN_DOCKER_REPO_USERNAME",
									ValueFrom: getSecretValue(kiteSecretsName, "docker-repo-username"),
								},
								{
									Name:      "TUFIN_DOCKER_REPO_PASSWORD",
									ValueFrom: getSecretValue(kiteSecretsName, "guru-api-key"),
								},
								{
									Name:      "API_KEY",
									ValueFrom: getSecretValue(kiteSecretsName, "guru-api-key"),
								},
							},
						},
					},
				},
			},
		},
	}
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
