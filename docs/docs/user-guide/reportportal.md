---
sidebar_position: 5
---


# ReportPortal Integration
Browserkube can automatically report test execution logs to ReportPortal
#### Configuration
Apply ReportPortal configuration in Browserkube namespace as Kubernetes secret:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: rp-secret # can be any orbitary string
  namespace: browserkube
  labels:
    io.browserkube.rp-project: "report-portal-project-name" # the label is important. This is how Browserkube understand that the secret belongs to ReportPortal Integration
type: Opaque
stringData:
  host: https://reportportal.epam.com # ReportPortal host
  authToken: "auth-token" # ReportPortal auth token
```
You are free to add as many secrets with project configuration as you want. 