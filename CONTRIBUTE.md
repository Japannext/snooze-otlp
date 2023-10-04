# Setup the golang environment

With gvm:
```bash
gvm use go1.20
```

Making sure you can build:
```bash
go build
```

Running the freshly built executable:
```bash
./snooze-otlp run
```

# Deploying in a local kubernetes cluster

If you have a local docker repo and kubernetes cluster, you
can prepare your environment like so:

Create a local `.make.env`
```bash
LOCAL_REPO=nexus.example.com:8443/snooze-otlp
# Change by your system custom CA you wan during the docker build
CA_BUNDLE=/usr/local/share/certificates/
```

> This will be used to upload your freshly built image.

Create a `.helmfile.yaml`:
```yaml
---
environments:
  default:
    # Replace it with the context of your kubernetes cluster
    kubeContext: dev

releases:
- name: snooze-otlp
  # Replace it with the namespace you will use for development
  namespace: myns
  chart: ./charts/snooze-otlp
  recreatePods: true
  values:
  - image:
      repo: nexus.example.com:8443/snooze-otlp
      tag: develop
      pullPolicy: Always
    replicas: 1
    logLevel: debug
    snooze:
      url: https://snooze.dev.example.com
      # Example of the CA bundle if using trust-manager.
      # You can create any configMap with your custom CA
      # And link it here.
      caConfigMap: ca-bundle
```

If you're working in a air-gap environment, make sure to configure
the `HTTPS_PROXY`/`NO_PROXY`:
```bash
export HTTPS_PROXY=proxy.example.com
export NO_PROXY=k8s.example.com
```

And if it's a TLS proxy:
```bash
make setup
```

> This will pull your current system CA to `./.ca-bundle/`, which docker will use.

Then, to build, package, and deploy:
```bash
make develop
```
