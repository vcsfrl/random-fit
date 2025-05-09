#!/bin/bash
kubectl delete -f job.yaml;
kubectl delete -f job-config.yaml;
kubectl delete -f processed-volumne.yaml;
rm processed-volumne.yaml;
