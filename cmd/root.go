// Copyright © 2022 Shyam Jeedigunta <shyam123.jvs95@gmail.com>
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

package cmd

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog/v2"
)

const (
	KubeStress = "kube-stress"
)

var (
	rootCmd = &cobra.Command{
		Use:   KubeStress,
		Short: "Simple tool for generating stress on a Kubernetes cluster.",
	}
	kubeconfig string
)

func init() {
	rootCmd.Flags().SortFlags = false
	klog.InitFlags(nil)
	klog.SetOutput(os.Stdout)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "Absolute path to the kubeconfig file")
}

func Execute() {
	defer klog.Flush()
	if err := rootCmd.Execute(); err != nil {
		klog.Errorf("Error executing root command: %v", err)
		os.Exit(1)
	}
}
