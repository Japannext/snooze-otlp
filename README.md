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

## Binary

WIP

# opentelemetry-collector

Here is an example of opentelemetry-collector configuration:
```yaml
---
# WIP
```
