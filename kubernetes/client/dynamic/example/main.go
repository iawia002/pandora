package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"

	"github.com/iawia002/pandora/kubernetes/client"
	unstructuredutils "github.com/iawia002/pandora/kubernetes/unstructured"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "kubeconfig",
				Usage: "kube config file path",
			},
		},
		Action: func(c *cli.Context) error {
			config, err := client.BuildConfigFromFlags("", c.String("kubeconfig"), client.SetQPS(100, 200))
			if err != nil {
				return err
			}
			dynamicClient, err := dynamic.NewForConfig(config)
			if err != nil {
				return err
			}
			objects, err := dynamicClient.Resource(corev1.SchemeGroupVersion.WithResource("nodes")).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return err
			}
			// metadata map[resourceVersion:91508551]
			// kind NodeList
			// apiVersion v1
			for key, value := range objects.Object {
				fmt.Println(key, value)
			}
			for i := range objects.Items {
				object := &objects.Items[i]
				conditions, _, err := unstructured.NestedSlice(object.Object, "status", "conditions")
				if err != nil {
					return err
				}
				node := &corev1.Node{}
				if err = unstructuredutils.ConvertToTyped(object, node); err != nil {
					return err
				}
				fmt.Printf("node %s podCIDR: %s, conditions length: %d\n", object.GetName(), node.Spec.PodCIDR, len(conditions))
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
