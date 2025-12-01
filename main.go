package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func main() {
	testEnv := &envtest.Environment{}

	// start testEnv
	cfg, err := testEnv.Start()
	if err != nil {
		log.Fatalf("Failed to start testEnv: %v", err)
	}

	defer testEnv.Stop()

	// Write kubeconfig file
	kubeconfigPath := filepath.Join(os.TempDir(), "envtest-kubeconfig")
	kubeconfig := clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			"envtest": {
				Server:                   cfg.Host,
				CertificateAuthorityData: cfg.CAData,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"envtest": {
				Cluster:  "envtest",
				AuthInfo: "envtest",
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"envtest": {
				ClientCertificateData: cfg.CertData,
				ClientKeyData:         cfg.KeyData,
			},
		},
		CurrentContext: "envtest",
	}
	err = clientcmd.WriteToFile(kubeconfig, kubeconfigPath)
	if err != nil {
		log.Fatalf("Failed to write kubeconfig: %v", err)
	}

	defer os.Remove(kubeconfigPath)

	fmt.Printf("Kubeconfig written to: %s\n", kubeconfigPath)
	fmt.Printf("You can use it with: export KUBECONFIG=%s\n", kubeconfigPath)
	fmt.Printf("Waiting for signal (Ctrl-C or kill %d) to stop...", os.Getpid())

	// Wait for signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\nReceived signal, shutting down...")
}
