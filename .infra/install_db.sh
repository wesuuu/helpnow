#!/bin/bash
set -e

NAMESPACE="helpnow"
CLUSTER_NAME="helpnow-postgres"

echo "Installing CloudNativePG Operator..."
kubectl apply -f https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.22/releases/cnpg-1.22.1.yaml
echo "Waiting for operator to be ready..."
kubectl wait --for=condition=Ready pod -l app.kubernetes.io/name=cloudnative-pg -n cnpg-system --timeout=60s

echo "Cleaning up previous installation..."
kubectl delete namespace $NAMESPACE --ignore-not-found=true
echo "Waiting for namespace deletion..."
kubectl wait --for=delete namespace/$NAMESPACE --timeout=60s || true

echo "Setting up helpnow namespace..."
kubectl create namespace $NAMESPACE

# echo "Applying Kustomize configuration..."
# # We run from project root, so -k .infra
kubectl apply -k .infra

echo "Waiting for Cluster to be ready..."
# Wait for 120s for primary
echo "Waiting for primary pod..."
kubectl wait --for=condition=Ready pod -l cnpg.io/cluster=$CLUSTER_NAME,role=primary -n $NAMESPACE --timeout=120s

echo "CloudNativePG Cluster 'helpnow-postgres' is ready!"
echo "Retrieving LoadBalancer IP..."
# Attempt to get IP or Hostname. For local k3s/docker, it might be localhost or pending.
SVC_IP=$(kubectl get svc helpnow-postgres-lb -n $NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}')

if [ -z "$SVC_IP" ]; then
    echo "Service IP not yet assigned (LoadBalancer pending?)."
    echo "If running locally (e.g. k3s/dockerd), port 5555 should be accessible on node IP."
else
    echo "Service Accessible at: $SVC_IP:5555"
fi

echo "Verifying schema execution..."
# Check for 'routines' table existence
kubectl exec -n $NAMESPACE $CLUSTER_NAME-1 -- psql -U postgres -d helpnow -c "\dt routines" || echo "Warning: Verification query failed."

echo "Installation Complete."
