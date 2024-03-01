package main

import (
	"context"
	"fmt"
	"os"

	"github.com/iawia002/lia/kubernetes/client"
	genericclient "github.com/iawia002/lia/kubernetes/client/generic"
	"github.com/urfave/cli/v2"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
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
			client, err := genericclient.NewClient(config,
				genericclient.WithContext(context.TODO()),
				genericclient.WithByObject(map[runtimeclient.Object]cache.ByObject{
					// store only metadata for Node objects
					&corev1.Node{}: {
						Transform: func(obj interface{}) (interface{}, error) {
							n, ok := obj.(*corev1.Node)
							if !ok {
								return obj, nil
							}
							n.ObjectMeta.ManagedFields = nil
							n.Spec = corev1.NodeSpec{}
							n.Status = corev1.NodeStatus{}
							return n, nil
						},
					},
				}),
			)
			if err != nil {
				return err
			}
			nodes := &corev1.NodeList{}
			if err = client.List(context.TODO(), nodes, runtimeclient.MatchingLabels{
				corev1.LabelOSStable: "linux",
			}); err != nil {
				return err
			}
			for _, node := range nodes.Items {
				fmt.Println(node)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
