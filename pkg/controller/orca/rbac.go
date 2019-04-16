package orca

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetServiceAccount(name string, namespace string, labels map[string]string) *corev1.ServiceAccount {

	return &corev1.ServiceAccount{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
	}
}

func GetClusterRole(name string, labels map[string]string, rules ...rbacv1.PolicyRule) *rbacv1.ClusterRole {

	return &rbacv1.ClusterRole{
		ObjectMeta: v1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Rules: rules,
	}
}

func GetPolicyRule(verbs []string, apiGroups []string, resources []string) rbacv1.PolicyRule {

	return rbacv1.PolicyRule{
		Verbs:     verbs,
		APIGroups: apiGroups,
		Resources: resources,
	}
}

func GetClusterRoleBindig(name string, serviceAccount *corev1.ServiceAccount, clusterRole *rbacv1.ClusterRole) *rbacv1.ClusterRoleBinding {

	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: serviceAccount.Namespace,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      serviceAccount.Name,
			Namespace: serviceAccount.Namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			APIGroup: clusterRole.APIVersion,
			Kind:     "ClusterRole",
			Name:     clusterRole.Name,
		},
	}
}
