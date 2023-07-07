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
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	foov1alpha1 "github.com/iawia002/pandora/kubernetes/apis/example/apis/foo/v1alpha1"
	controllerruntime "github.com/iawia002/pandora/kubernetes/controller/controller-runtime"
	samplecontroller "github.com/iawia002/pandora/kubernetes/controller/sample-controller"
	controllerwebhook "github.com/iawia002/pandora/kubernetes/controller/webhook"
	"github.com/iawia002/pandora/kubernetes/scheme"
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
		MetricsBindAddress:      ":8080",
		HealthProbeBindAddress:  ":8081",
		Scheme:                  scheme.Scheme,
		LeaderElection:          true,
		LeaderElectionNamespace: metav1.NamespaceSystem,
		LeaderElectionID:        "sample-controller-manager-leader-election",
		Logger:                  klog.NewKlogr(),
		SyncPeriod:              pointer.Duration(time.Hour * 1),
		// webhook config
		Port:    8443,
		CertDir: "./certs",
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

	nodeController, err := samplecontroller.NewController(kubeClient, informer.Core().V1().Nodes())
	if err != nil {
		return err
	}
	if err = mgr.Add(nodeController); err != nil {
		return err
	}

	// Register webhooks
	if err = (&foov1alpha1.Foo{}).SetupWebhookWithManager(mgr); err != nil {
		return err
	}
	mgr.GetWebhookServer().Register("/mutate-v1-pod", &webhook.Admission{Handler: &controllerwebhook.PodAnnotator{Client: mgr.GetClient()}})

	ctx := signals.SetupSignalHandler()
	informer.Start(ctx.Done())
	return mgr.Start(ctx)
}
