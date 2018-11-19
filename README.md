# nfs_exporter
NFS exporter for Prometheus

## Installation

```
go get -u -v github.com/aixeshunter/nfs_exporter

./${GOPATH}/bin/nfs_exporter --${flags} ...
```


## Usage of `nfs_exporter`

| Option                    | Default             | Description
| ------------------------- | ------------------- | -----------------
| -h, --help                | -                   | Displays usage.
| --web.listen-address      | `:9689`             | The address to listen on for HTTP requests.
| --web.metrics-path        | `/metrics`          | URL Endpoint for metrics
| --nfs.storage-path        | `/opt/nfs`          | The nfs storage mount path
| --nfs.address             | `127.0.0.1`         | The nfs server IP address
| --nfs.executable-path     | `/usr/sbin/showmount` | Path to showmount executable.
| --log.format              | `logger:stderr`     | Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
| --log.level               | `info`              | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
| --version                 | -                   | Prints version information


## Make
```
promu:  make prometheus library
build:  Go build
docker: build and run in docker container
gotest: run go tests and reformats
format: formatting code
vet:    vetting code
```

**build**: runs go build for nfs_exporter

**docker**: runs docker build and copy new built nfs_exporter

**gotest**: runs *vet* and *fmt* go tools


## Metrics

### Command: `showmount -e ${NFS_SERVER_IP}`

| Name          | type     | impl. state |
| ------------  | -------- | ------------|
| up            | Gauge    | implemented |


```sh
# TYPE nfs_up gauge
# HELP nfs_up Was the last query of NFS successful.
nfs_up{mount_path=" /mnt",nfs_address="192.168.0.2"} 0
nfs_up{mount_path=" /opt",nfs_address="192.168.0.3"} 0
nfs_up{mount_path="/opt/nfs",nfs_address="192.168.0.3"} 1
```

## Example in Prometheus of Kubernetes Cluster

[yaml file](prometheus/manifests)

### points

```yaml
      containers:
      - name: nfs-exporter
      
        image: aixeshunter/nfs_exporter:v1.0
        args:
        - "--nfs.storage-path=/opt/nfs1, /opt/nfs2, /opt/nfs3"   # NFS storage mount path
        - "--nfs.address=192.168.0.2"                            # NFS server IP address
        ports:
        - name: http-metrics
          containerPort: 9689
```