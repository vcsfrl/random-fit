#!/bin/bash
rm processed-volumne.yaml
kubectl delete -f job.yaml
kubectl delete -f job-config.yaml
kubectl delete -f processed-volumne.yaml
