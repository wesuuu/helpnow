#!/bin/bash
set -e

# Add Altinity Repo if not exists
if ! helm repo list | grep -q "altinity"; then
    helm repo add altinity https://altinity.github.io/helm-charts
fi
helm repo update

# Load env vars
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs)
fi

# Create Kubernetes Secret for ClickHouse
echo "Creating ClickHouse Credentials Secret..."
kubectl create secret generic clickhouse-credentials \
  --from-literal=username=$CLICKHOUSE_USER \
  --from-literal=password=$CLICKHOUSE_PASSWORD \
  --namespace helpnow \
  --dry-run=client -o yaml | kubectl apply -f -

# Install ClickHouse (Manual StatefulSet)
echo "Installing ClickHouse StatefulSet..."
kubectl apply -f .infra/clickhouse/clickhouse-statefulset.yaml

# Apply custom LoadBalancer service
echo "Applying custom LoadBalancer service..."
kubectl apply -f .infra/clickhouse/service.yaml
