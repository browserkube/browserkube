---
sidebar_position: 1
---

# Download Files

### How-to access downloaded files
Downloaded files are available when session is open, e.g. browser is not quit.
In order to download the file, send HTTP GET to the following endpoint:
```
{{BROWSERKUBE_URL}}/session/{{SESSION_ID}}/browserkube/downloads/{{FILE_NAME}}
```
for instance using curl,
```sh
curl -s -O myfile.txt http://kubernetes.docker.local/session/session-id/browserkube/downloads/myfile.txt
```

