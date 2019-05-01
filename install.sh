#!/bin/bash

args=""
ns="${1:-default}"
if [[ "$1" ]]; then
    args="-n ${1}"
    kubectl create ns "$1"
else
    args="-n default"
fi

kubectl ${args} apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
#kubectl ${args} apply -f deploy/role_binding.yaml
roles="kite conntrack orca-operator"
for role in $roles; do
    kubectl create clusterrolebinding "$role" --clusterrole "$role" --serviceaccount "${ns}:${role}"

    if [[ "$?" != "0" ]]; then
        kubectl delete clusterrolebinding "$role"
        kubectl create clusterrolebinding "$role" --clusterrole "$role" --serviceaccount "${ns}:${role}"
    fi
done

kubectl apply -f deploy/crds/tufin_v1alpha1_orca_crd.yaml
kubectl apply -f deploy/crds/tufin_v1_policy_crd.yaml
kubectl ${args} apply -f deploy/operator.yaml
