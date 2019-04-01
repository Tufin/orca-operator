# orca-operator
Kubernetes Operator for deploying Orca Agent (Kite etc.)

# After changing kite_types.go, run:
`operator-sdk generate k8s`

# Build
```
operator-sdk build tufin/orca-operator
docker push tufin/orca-operator
```
