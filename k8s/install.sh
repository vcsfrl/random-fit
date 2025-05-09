#!/bin/bash
set -a && source ./../.env && set +a
envsubst < volumne.yaml > processed-volumne.yaml
kubectl apply -f job-config.yaml
kubectl apply -f processed-volumne.yaml
kubectl apply -f job.yaml
