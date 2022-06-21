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
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	"github.com/shyamjvs/kube-stress/pkg/client"
	"github.com/shyamjvs/kube-stress/pkg/util"
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
			if err := createCommand(); err != nil {
				klog.Errorf("Error executing create command: %v", err)
				os.Exit(1)
			}
		},
	}
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createConfig.namespace, "namespace", KubeStress, "Namespace name where the test objects will be created")
	createCmd.Flags().StringVar(&createConfig.objectType, "object-type", "configmaps", "Type of objects to create (supported values are 'pods' and 'configmaps'")
	createCmd.Flags().IntVar(&createConfig.objectSize, "object-size-bytes", 40000, "Size of each object to be created (only used for 'configmap' object type)")
	createCmd.Flags().IntVar(&createConfig.objectCount, "object-count", 100, "Number of objects to create")
	createCmd.Flags().IntVar(&createConfig.numClients, "num-clients", 10, "Number of clients to use for spreading the create calls")
	createCmd.Flags().Float32Var(&createConfig.qps, "qps", 10.0, "QPS to use while creating the objects")
}

func createCommand() error {
	clients := client.CreateKubeClients(client.GetKubeConfig(kubeconfig), createConfig.numClients)

	// Setup signal handling for the process.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(context.Background())
	var once sync.Once
	defer once.Do(cancel)
	go func() {
		for {
			select {
			case sig := <-sigs:
				klog.V(1).Infof("Received stop signal: %v", sig)
				once.Do(cancel)
			case <-ctx.Done():
				klog.V(1).Info("Cancelled context and exiting program")
				return
			}
		}
	}()

	klog.V(1).Infof("Creating %v objects of type '%v' (%v bytes each) in namespace '%v' using %v clients and QPS = %v",
		createConfig.objectCount,
		createConfig.objectType,
		createConfig.objectSize,
		createConfig.namespace,
		createConfig.numClients,
		createConfig.qps)
	createObjects(ctx, clients)
	return nil
}

func createObjects(ctx context.Context, clients []*kubernetes.Clientset) {
	ticker := time.NewTicker(time.Second / time.Duration(createConfig.qps))
	defer ticker.Stop()

	var wg sync.WaitGroup
	var numObjectsCreated uint32
	for i := 0; atomic.LoadUint32(&numObjectsCreated) < uint32(createConfig.objectCount); i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			client := clients[i%len(clients)]
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := createObject(ctx, client); err == nil {
					atomic.AddUint32(&numObjectsCreated, 1)
				}
			}()
		}
	}

	wg.Wait()
	klog.V(1).Infof("Successfully created %v objects", numObjectsCreated)
}

func createObject(ctx context.Context, client *kubernetes.Clientset) error {
	start := time.Now()
	objectName := "configmap-" + uuid.Must(uuid.NewRandom()).String()
	// TODO: Implement other object-types below.
	configmap := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: objectName,
			// Below label helps selectively list/delete objects created by kube-stress.
			Labels: map[string]string{
				KubeStress: objectName,
			},
		},
		Data: map[string]string{objectName: util.RandomString(createConfig.objectSize)},
	}

	_, err := client.CoreV1().ConfigMaps(createConfig.namespace).Create(ctx, configmap, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Failed to create object: %v", err)
		return err
	}

	klog.V(2).Infof("Created object %v successfully (took %v)", objectName, time.Since(start))
	return nil
}
