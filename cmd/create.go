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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

type CreateConfig struct {
	namespace   string
	objectType  string
	objectSize  int
	objectCount int
	numClients  int
	qps         float32
}

var (
	createConfig *CreateConfig
	createCmd    *cobra.Command
)

func init() {
	createConfig = &CreateConfig{}
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create objects of a given type in the cluster",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runCommand(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		},
	}
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createConfig.namespace, "namespace", "kube-stress", "Namespace name where the test objects will be created")
	createCmd.Flags().StringVar(&createConfig.objectType, "object-type", "configmap", "Type of objects to create (supported values are 'pod' and 'configmap'")
	createCmd.Flags().IntVar(&createConfig.objectSize, "object-size-bytes", 40000, "Size of each object to be created (only used for 'configmap' object type)")
	createCmd.Flags().IntVar(&createConfig.objectCount, "object-count", 2000, "Number of objects to create")
	createCmd.Flags().IntVar(&createConfig.numClients, "num-clients", 10, "Number of clients to use for load-balancing the create calls")
	createCmd.Flags().Float32Var(&createConfig.qps, "qps", 10.0, "QPS to use while creating the objects")
}

func runCommand() error {
	klog.V(1).Infof("Creating %v objects of type '%v' (%v bytes each) in namespace '%v' using %v clients and total QPS = %v",
		createConfig.objectCount,
		createConfig.objectType,
		createConfig.objectSize,
		createConfig.namespace,
		createConfig.numClients,
		createConfig.qps)
	return nil
}
