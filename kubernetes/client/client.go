package client

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// DefaultQPS is the default QPS value.
	DefaultQPS = 50
	// DefaultBurst is the default Burst value.
	DefaultBurst = 100
)

// SetQPS sets the QPS and Burst.
func SetQPS(qps float32, burst int) func(c *rest.Config) {
	return func(c *rest.Config) {
		c.QPS = qps
		c.Burst = burst
	}
}

// BuildConfigFromFlags builds rest configs from a master url or a kube config filepath.
func BuildConfigFromFlags(masterURL, kubeConfigPath string, options ...func(c *rest.Config)) (*rest.Config, error) {
	c, err := clientcmd.BuildConfigFromFlags(masterURL, kubeConfigPath)
	if err != nil {
		return nil, err
	}

	for _, opt := range options {
		opt(c)
	}
	if c.QPS == 0 {
		c.QPS = DefaultQPS
	}
	if c.Burst == 0 {
		c.Burst = DefaultBurst
	}

	return c, nil
}
