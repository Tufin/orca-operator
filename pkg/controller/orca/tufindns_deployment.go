package orca

import (
	"strings"

	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"fmt"
)

func getTufinDNSDeployment(cr *appv1alpha1.Orca, dnsPort int32) *appsv1.Deployment {

	labels := GetLabels("k8s-app=" + kubeDNS)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      kubeDNS,
			Namespace: kubeSystem,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						GetConfigMapVolume("tufin-dns-config", "tufindns", corev1.KeyToPath{
							Key:  "Corefile",
							Path: "Corefile",
						}),
					},
					Containers: []corev1.Container{
						{
							Name:            tufinDNS,
							Image:           cr.Spec.Images["dns"],
							ImagePullPolicy: corev1.PullAlways,
							VolumeMounts: []corev1.VolumeMount{
								{Name: "tufin-dns-config", MountPath: "/etc/coredns"},
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: dnsPort, Protocol: corev1.ProtocolTCP, Name: "dns"},
								{ContainerPort: 54, Protocol: corev1.ProtocolTCP, Name: "dns-tcp"},
								{ContainerPort: 9154, Protocol: corev1.ProtocolTCP, Name: "metrics"},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "TUFIN_GRPC_DISCOVERY_URL",
									Value: fmt.Sprintf("%s.%s:81", kite, cr.Spec.Namespace),
								},
								{
									Name:  "TUFIN_FALLTHROUGH_DOMAINS",
									Value: cr.Spec.EndPoints["guru"],
								},
								{
									Name:  "IGNORED_CONFIGMAPS",
									Value: strings.Join(cr.Spec.IngnoredConfigMaps, ","),
								},
							},
						},
					},
				},
			},
		},
	}
}

func getTufinDNSConfigMap(cr *appv1alpha1.Orca, dnsPort int32) *corev1.ConfigMap {

	kubePort := 53

	if dnsPort == 53 {
		kubePort = 54
	}

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{APIVersion: "v1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tufinDNS,
			Namespace: kubeSystem,
		},
		Data: map[string]string{
			"Corefile": fmt.Sprintf(
				`
internal:%d {
  forward . 0.0.0.0:%d
}
.:%d {
  health :8091
  whitelist cluster.local {
    pods verified
  }
  proxy . 0.0.0.0:%d
}`, dnsPort, kubePort, dnsPort, kubePort),
		},
	}
}
