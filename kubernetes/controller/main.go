package main

import (
	"os"
	"time"

	"github.com/iawia002/lia/kubernetes/client"
	"github.com/urfave/cli/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	controllerruntime "github.com/iawia002/pandora/kubernetes/controller/controller-runtime"
	samplecontroller "github.com/iawia002/pandora/kubernetes/controller/sample-controller"
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
			config, err := client.BuildConfigFromFlags("", c.String("kubeconfig"))
			if err != nil {
				return err
			}
			return run(config)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func run(config *rest.Config) error {
	mgr, err := manager.New(config, manager.Options{
		LeaderElection:          true,
		LeaderElectionNamespace: metav1.NamespaceSystem,
		LeaderElectionID:        "sample-controller-manager-leader-election",
		Logger:                  klog.NewKlogr(),
		SyncPeriod:              pointer.Duration(time.Hour * 1),
		MetricsBindAddress:      ":8080",
		HealthProbeBindAddress:  ":8081",
	})
	if err != nil {
		return err
	}

	if err = mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return err
	}
	if err = mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return err
	}

	nodeReconciler := &controllerruntime.NodeReconciler{}
	if err = nodeReconciler.SetupWithManager(mgr); err != nil {
		return err
	}
	// Register timed tasker
	if err = mgr.Add(nodeReconciler); err != nil {
		return err
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	informer := informers.NewSharedInformerFactory(kubeClient, time.Second*30)
	nodeController := samplecontroller.NewController(informer.Core().V1().Nodes())
	if err = mgr.Add(nodeController); err != nil {
		return err
	}

	ctx := signals.SetupSignalHandler()
	informer.Start(ctx.Done())
	return mgr.Start(ctx)
}
