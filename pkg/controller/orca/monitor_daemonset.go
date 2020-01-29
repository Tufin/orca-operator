package orca

import (
	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

func getMonitorDaemonset(cr *appv1alpha1.Orca) *appsv1.DaemonSet {

	labels := GetLabels(app + "=" + monitor)

	if monitorImg := os.Getenv("RELATED_IMAGE_TUFIN_MONITOR"); monitorImg != "" {
		cr.Spec.Images[monitor] = monitorImg
	}

	return &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      monitor,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector:        GetLabelSelector(labels),
			MinReadySeconds: 5,
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: "RollingUpdate",
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      cr.Name,
					Namespace: cr.Namespace,
					Labels:    labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: monitor,
					DNSPolicy:          corev1.DNSClusterFirstWithHostNet,
					HostNetwork:        true,
					HostPID:            true,
					Tolerations: []corev1.Toleration{
						{Effect: corev1.TaintEffectNoSchedule, Operator: corev1.TolerationOpExists},
					},
					Volumes: []corev1.Volume{
						GetHostVolume(dockerSocketVolumeName, dockerSocketVolumePath, corev1.HostPathUnset),
						GetHostVolume(scopeKernelDebugVolumeName, scopeKernelDebugVolumePath, corev1.HostPathDirectoryOrCreate),
					},
					Containers: []corev1.Container{
						{
							Name:  monitor,
							Image: cr.Spec.Images[monitor],
							VolumeMounts: []corev1.VolumeMount{
								{Name: dockerSocketVolumeName, MountPath: dockerSocketVolumePath},
								{Name: scopeKernelDebugVolumeName, MountPath: scopeKernelDebugVolumePath},
							},
							Args: []string{
								"kite." + cr.Namespace + ":80",
							},
							Command: []string{"/home/tufin/monitor"},
							SecurityContext: &corev1.SecurityContext{
								Privileged: GetBoolRef(true),
							},
							Env: []corev1.EnvVar{
								{
									Name: "KUBERNETES_NODENAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											APIVersion: "v1",
											FieldPath:  "spec.nodeName",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
