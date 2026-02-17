## How to configure distributed tracing

1. Enable telemetry.providerEnabled, grafana.enable, tempo.enable in values.yaml

2.  Update dependencies 
```console
helm dependencies update ./helm/charts/browserkube
```

3. Authorization data for grafana:

- login: admin
- password: password