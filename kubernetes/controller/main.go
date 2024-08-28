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
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

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
	logger := klog.NewKlogr()
	log.SetLogger(logger)

	mgr, err := manager.New(config, manager.Options{
		Scheme:                  scheme.Scheme,
		LeaderElection:          true,
		LeaderElectionNamespace: metav1.NamespaceSystem,
		LeaderElectionID:        "sample-controller-manager-leader-election",
		Logger:                  logger,
		HealthProbeBindAddress:  ":8081",
		Metrics: metricsserver.Options{
			BindAddress: ":8080",
		},
		Cache: cache.Options{
			SyncPeriod: ptr.To(time.Hour * 1),
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Port:    8443,
			CertDir: "./certs",
		}),
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

	if err = (&controllerruntime.NodeReconciler{}).SetupWithManager(mgr); err != nil {
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
	if err = controllerwebhook.SetupFooWebhookWithManager(mgr); err != nil {
		return err
	}
	if err = controllerwebhook.SetupPodWebhookWithManager(mgr); err != nil {
		return err
	}

	ctx := signals.SetupSignalHandler()
	informer.Start(ctx.Done())
	return mgr.Start(ctx)
}
