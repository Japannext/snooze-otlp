# snooze-otlp

A Snooze opentelemetry log input plugin.

Listen on gRPC for opentelemetry /v1/logs for log inputs,
and forward them to a snooze instance in the correct format.

# Deployment

## Helm chart

```bash
helm repo add jnx japannext.github.com/helm-charts
helm install jnx/snooze-otlp -f values.yaml
```

Example of a basic configuration:
```yaml
# values.yaml
snooze:
  url: https://snooze.example.com
  caConfigMap: ca-bundle  # from trust-manager for instance
```

## Docker

```bash
docker run ghcr.io/japannext/snooze-otlp:latest \
    --env-file env.list \
    -v /etc/pki/tls/cert.pem:/tls/ca.crt:ro
```

```bash
# env.list
SNOOZE_OTLP_SNOOZE_URL=https://snooze.example.com
SNOOZE_OTLP_CA_PATH=/tls/ca.crt
```

# opentelemetry-collector

Here is an example of opentelemetry-collector configuration.
This create an instance of otelcol that listen on OTLP log protocol, append kubernetes pod metadata,
and forward it to snooze (via snooze-otlp).
```yaml
---
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: ${MY_POD_IP}:4317
      http:
        endpoint: ${MY_POD_IP}:4318
processors:
  k8sattributes:
    extract:
      metadata:
      - k8s.namespace.name
      - k8s.deployment.name
      - k8s.statefulset.name
      - k8s.daemonset.name
      - k8s.cronjob.name
      - k8s.job.name
    passthrough: false
    pod_association:
    - sources:
      - from: resource_attribute
        # Allow pods to send OTLP logs to this instance,
        # and have their metadata detected based on their
        # source IP.
        name: k8s.pod.ip
    - sources:
      - from: resource_attribute
        name: k8s.pod.uid
    - sources:
      - from: connection
exporters:
  otlp/snooze:
    compression: none
    # When snooze-otlp is deployed in the same namespace as otel-collector
    endpoint: snooze-otlp:4317
    tls:
      insecure: true
service:
  pipelines:
    logs:
      receivers:
      - otlp
      processors:
      - k8sattributes
      exporters:
      - otlp/snooze
```
