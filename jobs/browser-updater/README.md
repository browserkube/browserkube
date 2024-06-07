# Browserkube Browser Updater

## How it works
At the start of the every helm install(not upgrade) browser updater will update the browserkube browserset with the custom or default registry and cache the
said images in kubernetes by start/stopping them as a browser. its purpose is keep the browserset up to date.

There are two instances that runs this updater. One is at the installation of the browserkube, and one at every midnight as a cronjob

## How to add new browser to the updater
To add a new browser image chain to the updater, simply add the image you need to keep track of to the browserset.yaml file, updater will search the new versions of that image when its being run

## How to add custom registry
1-Crate secrets to the custom registries with .dockerconfig json data field encoded to base64(More info:https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/)
This secrets hold the details of your registry such as registry url, username, password etc. These registries HAS to support Docker Registry V2 API.\
To add a registry without auth leave the username, password etc fields as empty strings.\
You can add multiple secrets to pull different images but to do so create a separate secret for each registry
2-In the browserset yaml file add the registrySecret(string) field to each version you want to use custom registry. this field should be the name of the secret you created.\
Example:
```yaml
  webdriver:
    firefox:
      defaultVersion: "108.0"
      defaultPath: "/wd/hub"
      versions:
        "108.0":
          image: selenium/standalone-firefox:108.0
          port: "4444"
          provider: k8s
          registrySecret: yoursecret1 # Your first secret
    chrome:
      defaultVersion: "111.0"
      defaultPath: "/"
      versions:
        "108.0":
          image: selenium/standalone-chrome:108.0
          port: "4444"
          provider: k8s
          registrySecret: yoursecret2
        "109.0":
          image: selenium/standalone-chrome:109.0
          port: "4444"
          provider: k8s
          registrySecret: yoursecret2 # Your second secret
```
### Notes:
*  If you dont add any secrets, updater will use docker hub without auth as a default registry
*  If you're using a repo different than docker hub, make sure to include the full name of the repo (Ex: mcr.microsoft.com/playwright:v1.0 rather than playwright:v1.0)

