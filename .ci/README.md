## Infra README

### Install MicroK8S cluster
#### Install MicroK8s
```sh
sudo snap install microk8s --classic
```
```
You can either try again with sudo or add the user <user-name> to the 'microk8s' group:

    sudo usermod -a -G microk8s <user-name>
    sudo chown -f -R <user-name> ~/.kube
```

#### Enable Ingress. Make sure DNS is set (if needed)
```sh
microk8s.enable registry:size=30Gi 
microk8s.enable ingress
```

#### Configure Internal Registry if needs to be available by hostname
```sh
sudo mkdir -p /var/snap/microk8s/current/args/certs.d/<host-port>
sudo touch /var/snap/microk8s/current/args/certs.d/<host-port>/hosts.toml
```

```sh
sudo mkdir -p /var/snap/microk8s/current/args/certs.d/registry.container-registry.svc.cluster.local:5000
sudo touch /var/snap/microk8s/current/args/certs.d/registry.container-registry.svc.cluster.local:5000/hosts.toml
```
Content of hosts.toml:
```
# /var/snap/microk8s/current/args/certs.d/registry.container-registry.svc.cluster.local:5000/hosts.toml
server = "registry.container-registry.svc.cluster.local:5000"

[host."http://registry.container-registry.svc.cluster.local:5000"]
capabilities = ["pull", "resolve"]
```

### If internal DNS is enabled
```sh
sudo mkdir -p /var/snap/microk8s/current/args/certs.d/<<host-name>>:32000
sudo touch /var/snap/microk8s/current/args/certs.d/<<host-name>>:32000/hosts.toml
```
Content of hosts.toml:
```
# /var/snap/microk8s/current/args/certs.d/<<host-name>>:32000/hosts.toml
server = "<<host-name>>:32000"

[host."<<host-name>>:32000"]
capabilities = ["pull", "resolve"]
```

### Example of microk8s configuration
`/var/snap/microk8s/current/args/containerd-template.toml`
```toml
# Use config version 2 to enable new configuration fields.
version = 2
oom_score = 0

[grpc]
  uid = 0
  gid = 0
  max_recv_message_size = 16777216
  max_send_message_size = 16777216

[debug]
  address = ""
  uid = 0
  gid = 0

[metrics]
  address = "127.0.0.1:1338"
  grpc_histogram = false

[cgroup]
  path = ""


# The 'plugins."io.containerd.grpc.v1.cri"' table contains all of the server options.
[plugins."io.containerd.grpc.v1.cri"]

  stream_server_address = "127.0.0.1"
  stream_server_port = "0"
  enable_selinux = false
  sandbox_image = "k8s.gcr.io/pause:3.1"
  stats_collect_period = 10
  enable_tls_streaming = false
  max_container_log_line_size = 16384

  # 'plugins."io.containerd.grpc.v1.cri".containerd' contains config related to containerd
  [plugins."io.containerd.grpc.v1.cri".containerd]

    # snapshotter is the snapshotter used by containerd.
    snapshotter = "${SNAPSHOTTER}"

    # no_pivot disables pivot-root (linux only), required when running a container in a RamDisk with runc.
    # This only works for runtime type "io.containerd.runtime.v1.linux".
    no_pivot = false

    # default_runtime_name is the default runtime name to use.
    default_runtime_name = "${RUNTIME}"

    # 'plugins."io.containerd.grpc.v1.cri".containerd.runtimes' is a map from CRI RuntimeHandler strings, which specify types
    # of runtime configurations, to the matching configurations.
    # In this example, 'runc' is the RuntimeHandler string to match.
    [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.runc]
      # runtime_type is the runtime type to use in containerd e.g. io.containerd.runtime.v1.linux
      runtime_type = "${RUNTIME_TYPE}"

    [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.nvidia-container-runtime]
      # runtime_type is the runtime type to use in containerd e.g. io.containerd.runtime.v1.linux
      runtime_type = "${RUNTIME_TYPE}"

      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.nvidia-container-runtime.options]
        BinaryName = "nvidia-container-runtime"

   [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata]
      runtime_type = "io.containerd.kata.v2"
      [plugins."io.containerd.grpc.v1.cri".containerd.runtimes.kata.options]
        BinaryName = "kata-runtime"

  # 'plugins."io.containerd.grpc.v1.cri".cni' contains config related to cni
  [plugins."io.containerd.grpc.v1.cri".cni]
    # bin_dir is the directory in which the binaries for the plugin is kept.
    bin_dir = "${SNAP_DATA}/opt/cni/bin"

    # conf_dir is the directory in which the admin places a CNI conf.
    conf_dir = "${SNAP_DATA}/args/cni-network"

  # 'plugins."io.containerd.grpc.v1.cri".registry' contains config related to the registry
  [plugins."io.containerd.grpc.v1.cri".registry]
    config_path = "${SNAP_DATA}/args/certs.d"
    [plugins."io.containerd.grpc.v1.cri".registry.mirrors]
   	  [plugins."io.containerd.grpc.v1.cri".registry.mirrors."registry.container-registry.svc.cluster.local:5000"]
        endpoint = ["http://registry.container-registry.svc.cluster.local:5000"]
      [plugins."io.containerd.grpc.v1.cri".registry.mirrors."<<host-name>>:32000"]
        endpoint = ["<<host-name>>:32000"]

```

### Configuring Skaffold Profile
```sh
eval $(sops --decrypt epm.enc.env)
```


### Install Tekton
```sh
# replace DASHBOARD_URL with the hostname you want for your dashboard
# the hostname should be setup to point to your ingress controller
DASHBOARD_URL=dashboard.domain.tld
kubectl apply -n tekton-pipelines -f - <<EOF
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tekton-dashboard
  namespace: tekton-pipelines
spec:
  rules:
  - host: $DASHBOARD_URL
    http:
      paths:
      - pathType: ImplementationSpecific
        backend:
          service:
            name: tekton-dashboard
            port:
              number: 9097
EOF
```
### Configure Tekton
```sh
tkn hub install task git-clone
tkn hub install task kaniko
```