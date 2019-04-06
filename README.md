# orca-operator
Kubernetes Operator for deploying Orca Agent (Kite etc.)

# After changing kite_types.go, run:
```
operator-sdk generate k8s
```

# After changing the CRD:
```
kubectl create -f deploy/crds/app_v1alpha1_kite_crd.yaml
```

# Build
```
operator-sdk generate k8s
operator-sdk build tufin/orca-operator
docker push tufin/orca-operator
```

# Run
```
kubectl create -f deploy/operator.yaml
```

# Create the Kite CR - the default controller will watch for Kite objects and install the agent
```
kubectl create -f deploy/crds/app_v1alpha1_kite_cr.yaml
```

# Details
https://github.com/operator-framework/operator-sdk
