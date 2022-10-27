package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	corev1 "k8s.io/api/core/v1"

	"github.com/iawia002/pandora/kubernetes/client"
	genericclient "github.com/iawia002/pandora/kubernetes/client/generic"
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
			client, err := genericclient.NewClient(config, genericclient.WithContext(context.TODO()))
			if err != nil {
				return err
			}
			nodes := &corev1.NodeList{}
			if err = client.List(context.TODO(), nodes); err != nil {
				return err
			}
			for _, node := range nodes.Items {
				fmt.Println(node.Name)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
