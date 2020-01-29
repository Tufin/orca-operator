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
roles="kite monitor orca-operator"
for role in $roles; do
    kubectl create clusterrolebinding "$role" --clusterrole "$role" --serviceaccount "${ns}:${role}"

    if [[ "$?" != "0" ]]; then
        kubectl delete clusterrolebinding "$role"
        kubectl create clusterrolebinding "$role" --clusterrole "$role" --serviceaccount "${ns}:${role}"
    fi
done

kubectl apply -f deploy/crds/orca-operator-orca.crd.yaml
kubectl apply -f deploy/crds/orca-operator-policy.crd.yaml
kubectl ${args} apply -f deploy/operator.yaml
