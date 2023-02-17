// Copyright CERN.
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cernops/cvmfs-csi/internal/cvmfs/automount"
	"github.com/cernops/cvmfs-csi/internal/log"
	cvmfsversion "github.com/cernops/cvmfs-csi/internal/version"

	"k8s.io/klog/v2"
)

var (
	version = flag.Bool("version", false, "Print driver version and exit.")

	hasAlienCache = flag.Bool("has-alien-cache", false, "CVMFS client is using alien cache volume")

	automountDaemonUnmountAfterIdleSeconds = flag.Int("automount-unmount-timeout", 300, "number of seconds of idle time after which an autofs-managed CVMFS mount will be unmounted. '0' means never unmount, '-1' leaves automount default option.")
)

func main() {
	// Handle flags and initialize logging.

	klog.InitFlags(nil)
	if err := flag.Set("logtostderr", "true"); err != nil {
		klog.Exitf("failed to set logtostderr flag: %v", err)
	}
	flag.Parse()

	if *version {
		fmt.Println("automount-runner for CVMFS CSI plugin version", cvmfsversion.FullVersion())
		os.Exit(0)
	}

	// Initialize and run automount-runner.

	log.Infof("automount-runner for CVMFS CSI plugin version %s", cvmfsversion.FullVersion())
	log.Infof("Command line arguments %v", os.Args)

	err := automount.Init(&automount.Opts{
		UnmountTimeoutSeconds: *automountDaemonUnmountAfterIdleSeconds,
		HasAlienCache:         *hasAlienCache,
	})

	if err != nil {
		log.Fatalf("Failed to initialize automount-runner: %v", err)
	}

	if err = automount.RunBlocking(); err != nil {
		log.Fatalf("Failed to run automount-runner: %v", err)
	}

	os.Exit(0)
}