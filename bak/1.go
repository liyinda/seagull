/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/retry"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
        "encoding/json"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

    //创建deployment
    go createDeployment(deploymentsClient)

    //监听deployment
    //startWatchDeployment(deploymentsClient)

}

//监听Deployment变化
func startWatchDeployment(deploymentsClient v1beta1.DeploymentInterface) {
    w, _ := deploymentsClient.Watch(metav1.ListOptions{})
    for {
        select {
        case e, _ := <-w.ResultChan():
            fmt.Println(e.Type, e.Object)
        }
    }
}

//创建deployemnt，需要谨慎按照部署的k8s版本来使用api接口
func createDeployment(deploymentsClient appsv1.DeploymentInterface)  {
    var r apiv1.ResourceRequirements
    //资源分配会遇到无法设置值的问题，故采用json反解析
    j := `{"limits": {"cpu":"2000m", "memory": "1Gi"}, "requests": {"cpu":"2000m", "memory": "1Gi"}}`
    json.Unmarshal([]byte(j), &r)
    deployment := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name: "engine",
            Labels: map[string]string{
                "app": "engine",
            },
        },
        Spec: appsv1.DeploymentSpec{
            Replicas: int32Ptr2(1),
            Template: apiv1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app": "engine",
                    },
                },
                Spec: apiv1.PodSpec{
                    Containers: []apiv1.Container{
                        {   Name:               "engine",
                            Image:           "my.com/engine:v2",
                            Resources: r,
                        },
                    },
                },
            },
        },
    }

    fmt.Println("Creating deployment...")
    result, err := deploymentsClient.Create(deployment)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }

