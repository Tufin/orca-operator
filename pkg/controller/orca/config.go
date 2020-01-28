package orca

import (
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	app     = "app"
	name    = "name"
	kite    = "kite"
	monitor = "monitor"

	dockerSocketVolumeName     = "docker-socket"
	dockerSocketVolumePath     = "/var/run/docker.sock"
	scopeKernelDebugVolumeName = "sys-kernel-debug"
	scopeKernelDebugVolumePath = "/sys/kernel/debug"
)

func GetBoolRef(val bool) *bool {

	ret := val

	return &ret
}

func BoolToString(b bool) string {

	return strconv.FormatBool(b)
}

func GetLabels(labels ...string) map[string]string {

	var ret map[string]string
	ret = make(map[string]string)

	for _, label := range labels {

		tmp := strings.Split(label, "=")
		if len(tmp) == 2 {
			ret[tmp[0]] = tmp[1]
		}
	}

	return ret
}

func GetLabelSelector(labels map[string]string) *metav1.LabelSelector {

	return &metav1.LabelSelector{MatchLabels: labels}
}

func GetHostVolume(name string, path string, volType corev1.HostPathType) corev1.Volume {

	return corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{Path: path, Type: &volType},
		},
	}
}

func GetConfigMapVolume(volName string, configName string, items ...corev1.KeyToPath) corev1.Volume {

	return corev1.Volume{
		Name: volName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configName,
				},
				Items: items,
			},
		},
	}
}
