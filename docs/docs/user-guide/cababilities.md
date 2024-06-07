---
sidebar_position: 3
---

# Capabilities

### Browserkube-specific capabilities
Browserkube capabilities are specified under ```browserkube:options``` key, e.g.:
```go
var caps = selenium.Capabilities{
	"browserName": "chrome",
	"enableVNC":   true,
	"browserkube:options": map[string]interface{}{
		"name":   "test session name", // session name
	},
}
```

