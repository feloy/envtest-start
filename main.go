package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func main() {
	users := flag.Int("users", 0, "add users to the testenv cluster, named user1, user2, user3, etc.")
	flag.Parse()

	testEnv := &envtest.Environment{}

	// start testEnv
	cfg, err := testEnv.Start()
	if err != nil {
		log.Fatalf("Failed to start testEnv: %v", err)
	}

	defer testEnv.Stop()

	var kubeconfigPath string
	if len(flag.Args()) > 0 {
		kubeconfigPath = flag.Arg(0)
	} else {
		kubeconfigPath = filepath.Join(os.TempDir(), "envtest-kubeconfig")
	}
	// Write kubeconfig file
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

	if users != nil && *users > 0 {
		for i := 0; i < *users; i++ {
			userInfo := envtest.User{
				Name: fmt.Sprintf("user%d", i+1),
			}
			user, err := testEnv.AddUser(userInfo, nil)
			if err != nil {
				log.Fatalf("Failed to create user: %v", err)
			}

			kubeConfig, err := user.KubeConfig()
			if err != nil {
				log.Fatalf("Failed to create the testenv-admin user kubeconfig: %v", err)
			}
			userKubeconfigPath := filepath.Join(path.Dir(kubeconfigPath), fmt.Sprintf("user%d-kubeconfig", i+1))
			err = os.WriteFile(userKubeconfigPath, kubeConfig, 0644)
			if err != nil {
				log.Fatalf("Failed to write kubeconfig: %v", err)
			}
			fmt.Printf("Kubeconfig for user%d written to: %s\n", i+1, userKubeconfigPath)
			defer os.Remove(userKubeconfigPath)
		}
	}

	fmt.Printf("Waiting for signal (Ctrl-C or kill %d) to stop...", os.Getpid())

	// Wait for signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	fmt.Println("\nReceived signal, shutting down...")
}
