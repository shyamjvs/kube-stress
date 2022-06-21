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

package client

import (
	"math"
	"os"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/klog/v2"
)

// Get a kubeconfig object from the supplied file path.
func GetKubeConfig(kubeconfig string) *restclient.Config {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		klog.Errorf("Error reading the kubeconfig file: %v", err)
		os.Exit(1)
	}
	// Disable client-go rate-limiting, we'll manage the test throughput ourselves.
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(math.MaxFloat32, math.MaxInt)
	return config
}

// Create a given number of k8s clients using the provided kubeconfig.
func CreateKubeClients(config *restclient.Config, numClients int) []*kubernetes.Clientset {
	clients := make([]*kubernetes.Clientset, 0, numClients)
	for i := 0; i < numClients; i++ {
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			klog.Errorf("Error creating a k8s client: %v", err)
			os.Exit(1)
		}
		clients = append(clients, client)
	}
	return clients
}
