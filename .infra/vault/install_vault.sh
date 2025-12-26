#!/bin/bash

set -e

helm repo add hashicorp https://helm.releases.hashicorp.com
helm upgrade --install vault hashicorp/vault --namespace vault --create-namespace -f $(dirname "$0")/values.yaml