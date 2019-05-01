# orca-operator
Kubernetes Operator for deploying Orca Agent (Kite etc.)


# Installing the operator on non-Openshift (internal use)
Run the following command - you can choose the namespace to deploy the operator in:
```bash
> ./install [OPTIONAL_NAMESPACE]
```

# Generating Orca's agent secret
1. Get the agent's token (kite) from Orca's UI in the setting tab
2. Get the domain & projects names from the UI's header
3. Run the following: 
```bash
> ./generate_orca_secret.sh <DOMAIN> <PROJECT> <AGENT_TOKEN> [OPTIONAL NAMESPACE] > orca_secret.yaml
> kubectl create -f orca_secret.yaml
```

# Create the Orca CR
Fill in the Domain,Project & Namespace - The agent will be installed in the supplied Namespace.
The controller will watch this CRD and install the Orca agent components (kite, conntrack etc.):
```
kubectl create -f deploy/crds/tufin_v1alpha1_orca_cr.yaml
```

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


