package main

import (
	"strings"
	"testing"

	"github.com/prometheus/common/log"
)

const (
	ExecPath = "/usr/sbin/showmount"
)

type testCase struct {
	mountpath []string
	address   string
}

func TestNFSMountPath(t *testing.T) {
	var tests = []testCase{
		{
			mountpath: []string{"/opt/nfs1", "/opt/nfs2", "/opt/nfs3"},
			address:   "127.0.0.1",
		},
		{
			mountpath: []string{"/var/lib/nfs1", "/var/lib/nfs2", "/var/lib/nfs3"},
			address:   "192.168.0.3",
		},
	}

	for _, c := range tests {
		for _, path := range c.mountpath {
			out, exists := execCommand(ExecPath, c.address)
			if !exists {
				t.Error("NFS not exist in server", c.address)
			}

			for _, p := range strings.Split(string(out), "\n") {
				if strings.Split(p, " ")[0] == path {
					log.Infoln("Mount Path is matching NFS server.", path, c.address)
					break
				}
				t.Errorf("Path is not found in NFS server %s.", path)
			}
		}
	}
}
