package orca

import (
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"strconv"
)

const (
	app          = "app"
	name         = "name"
	kite         = "kite"
	conntrack    = "conntrack"
	orcaOperator = "orca-operator"

	dockerSocketVolumeName     = "docker-socket"
	dockerSocketVolumePath     = "/var/run/docker.sock"
	scopePluginsVolumeName     = "scope-plugins"
	scopePluginsVolumePath     = "/var/run/scope/plugins"
	scopeKernelDebugVolumeName = "sys-kernel-debug"
	scopeKernelDebugVolumePath = "/sys/kernel/debug"

	VerbAll    = "*"
	VerbGet    = "get"
	VerbWatch  = "watch"
	VerbList   = "list"
	VerbCreate = "create"
	VerbUpdate = "update"
	VerbDelete = "delete"

	ResourceNodes                  = "nodes"
	ResourceIngresses              = "ingresses"
	ResourceEndpoints              = "endpoints"
	ResourceConfigMaps             = "configmaps"
	ResourceNamespaces             = "namespaces"
	ResourceServices               = "services"
	ResourceDeployments            = "deployments"
	ResourceDaemonSets             = "daemonsets"
	ResourcePods                   = "pods"
	ResourceSecrets                = "secrets"
	ResourceStorageClasses         = "storageclasses"
	ResourcePersistentVolumes      = "persistentvolumes"
	ResourcePersistentVolumeClaims = "persistentvolumeclaims"
	ResourceCronJobs               = "cronjobs"
	ResourceJobs                   = "jobs"
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
