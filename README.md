# kube-annotate

[![Build Status](https://travis-ci.org/chickenzord/kube-annotate.svg?branch=master)](https://travis-ci.org/chickenzord/kube-annotate)
[![Go Report Card](https://goreportcard.com/badge/github.com/chickenzord/kube-annotate)](https://goreportcard.com/report/github.com/chickenzord/kube-annotate)
[![codecov](https://codecov.io/gh/chickenzord/kube-annotate/branch/master/graph/badge.svg)](https://codecov.io/gh/chickenzord/kube-annotate)
[![Automated Docker Build](https://img.shields.io/docker/automated/chickenzord/kube-annotate.svg)](https://hub.docker.com/r/chickenzord/kube-annotate/) 
[![Docker Pulls](https://img.shields.io/docker/pulls/chickenzord/kube-annotate.svg)](https://hub.docker.com/r/chickenzord/kube-annotate/)

Kubernetes mutating admission webhook to automatically annotate pods.

Features:
- Automatically annotate new pods with certain labels
- YAML-based configuration for multiple rules
- Built-in Prometheus metrics exporter

Configurations:

- LOG_FORMAT: json/text
- LOG_LEVEL: trace/debug/info/warning/error/fatal/panic
- RULES_FILE: path to `config.yaml`
- TLS_ENABLED: must be `true` when running inside Kubernetes cluster as admission controller
- TLS_CRT: path to certfile for TLS config 
- TLS_KEY: path to keyfile for TLS config

Rules config sample:

```yaml
# config.yaml
- selector:
    app: http-service
  annotations:
    log.config.scalyr.com/include: true
- selector:
    app: postgresql
  annotations:
    log.config.scalyr.com/include: false
```

Setup:

1. Make sure the cluster support admission controller (at least Kubernetes 1.9)
2. Prepare TLS certificate (see Medium post below, you need cluster-admin permission)
3. Create kubernetes resources (see `examples` directory, please read the comments especially about CA bundle and certificates)
4. Label the namespace you want to enable (`kubectl label namespace ${namespace} kube-annotate=enabled`)

---

TODO:
- ~~bind internal endpoints (health, metrics) to separate port~~
- proper request/response logging
- ~~prometheus exporter~~
- helm chart for easier setup

---
References: 
- https://medium.com/ibm-cloud/diving-into-kubernetes-mutatingadmissionwebhook-6ef3c5695f74
- https://github.com/morvencao/kube-mutating-webhook-tutorial