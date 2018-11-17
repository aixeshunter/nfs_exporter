// Copyright 2018 Aixes Hunter
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// NFS exporter, exports metrics from Linux commandline tool like showmount.
package main

import (
	"net/http"

	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	nfsCmd       = "/usr/sbin/showmount"
	namespace    = "nfs"
	nfsMountPath = "/opt/nfs"
	nfsAddress   = "127.0.0.1"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last query of NFS successful.",
		[]string{"mount_path", "nfs_address"}, nil,
	)
)

// Exporter holds name, path and volumes to be monitored
type Exporter struct {
	hostname  string
	execpath  string
	address   string
	mountpath []string
}

// Describe all the metrics exported by NFS exporter. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

// Collect collects all the metrics
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Collect metrics from volume info
	out, exists := execCommand(e.execpath, e.address)

	for _, path := range e.mountpath {
		flag := 0.0

		if !exists {
			log.Fatalf("Get NFS storage path failed caused by: %s", string(out))
		} else {
			log.Infoln("Getting showmount result succeed.")

			for _, p := range strings.Split(string(out), "\n") {
				if strings.Split(p, " ")[0] == path {
					log.Infoln("Mount Path is matching NFS server.")
					flag = 1.0
					break
				}
			}
		}

		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, flag,
			path, e.address,
		)
	}
}

// NewExporter initialises exporter
func NewExporter(hostname, nfsExcPath, nfsPath, nfsAddress string) (*Exporter, error) {
	if len(nfsExcPath) < 1 {
		log.Fatalf("NFS executable path is wrong: %v", nfsExcPath)
	}
	volumes := strings.Split(nfsPath, ",")
	if len(volumes) < 1 {
		log.Warnf("No NFS storage mount path given. Proceeding without path information. Path: %v", nfsPath)
	}

	return &Exporter{
		hostname:  hostname,
		execpath:  nfsExcPath,
		address:   nfsAddress,
		mountpath: volumes,
	}, nil
}

func init() {
	prometheus.MustRegister(version.NewCollector("nfs_exporter"))
}

func main() {
	// commandline arguments
	var (
		metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
		listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9689").String()
		nfsExcPath    = kingpin.Flag("nfs.executable-path", "Path to nfs executable.").Default(nfsCmd).String()
		nfsPath       = kingpin.Flag("nfs.storage-path", "Path to nfs storage volume.").Default(nfsMountPath).String()
		nfsAddress    = kingpin.Flag("nfs.address", "IP address to nfs storage cluster.").Default("127.0.0.1").String()
		num           int
	)

	log.AddFlags(kingpin.CommandLine)
	kingpin.Version(version.Print("nfs_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Infoln("Starting nfs_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("While trying to get Hostname error happened: %v", err)
	}
	exporter, err := NewExporter(hostname, *nfsExcPath, *nfsPath, *nfsAddress)
	if err != nil {
		log.Errorf("Creating new Exporter went wrong, ... \n%v", err)
	}
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		num, err = w.Write([]byte(`<html>
			<head><title>NFS Exporter v` + version.Version + `</title></head>
			<body>
			<h1>NFS Exporter v` + version.Version + `</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Fatal(num, err)
		}
	})

	log.Infoln("Listening on", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
