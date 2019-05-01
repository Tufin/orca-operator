# orca-operator
Kubernetes Operator for deploying Orca Agent (Kite etc.)

# Build
```
operator-sdk generate k8s # only needed after changing kite_types.go
operator-sdk build tufin/orca-operator --image-build-args "--build-arg version=${CIRCLE_BUILD_NUM} --build-arg release=${OPENSHIFT_RELEASE}"
docker push tufin/orca-operator
```

# Setup k8s RBAC 
```
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
```

# Create the CRD:
```
kubectl create -f deploy/crds/tufin_v1alpha1_orca_crd.yaml
```

# Run the operator
```
kubectl create -f deploy/operator.yaml
```

# Create the Kite CR
The controller will watch for these CRD and install the Orca agent components (kite, conntrack etc.):
```
kubectl create -f deploy/crds/tufin_v1alpha1_orca_cr.yaml
```

# Deploying on RH OpenShift
In order to deploy in OpenShift we need to add `SecurityContextConstraints` 
to 'the operator, kite & conntrack' service accounts as follows:

```
# for kite
oc adm policy add-scc-to-user hostaccess -z kite -n tufin-system
oc adm policy add-scc-to-user hostnetwork -z kite -n tufin-system
oc adm policy add-scc-to-user node-exporter -z kite -n tufin-system

# for conntrack
oc adm policy add-scc-to-user privileged -z conntrack -n tufin-system
```

# Details
https://github.com/operator-framework/operator-sdk


