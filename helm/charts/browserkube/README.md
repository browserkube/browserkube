## How to configure distributed tracing

1. Enable telemetry.providerEnabled, grafana.enable, tempo.enable in values.yaml

2. Choose telemetry.providerType in browserkube-backend-deployment.yaml:

- zipkin
- otlptracehttp

```yaml
- name: TELEMETRY_PROVIDER_TYPE
  value: {{ .Values.telemetry.providerType.zipkin }}
```

3.  Update dependencies 
```console
helm dependencies update ./helm/charts/browserkube
```

4. Authorization data for grafana:

- login: admin
- password: password