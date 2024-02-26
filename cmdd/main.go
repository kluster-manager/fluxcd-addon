package main

import (
	"context"
	"fmt"
	"github.com/kluster-manager/fluxcd-addon/apis/fluxcd/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2/klogr"
	"open-cluster-management.io/addon-framework/pkg/addonfactory"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/yaml"
)

func NewClient() (client.Client, error) {
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	// NOTE: Register KubeDB api types
	_ = v1alpha1.AddToScheme(scheme)

	ctrl.SetLogger(klogr.New())
	cfg := ctrl.GetConfigOrDie()
	cfg.QPS = 100
	cfg.Burst = 100

	hc, err := rest.HTTPClientFor(cfg)
	if err != nil {
		return nil, err
	}
	mapper, err := apiutil.NewDynamicRESTMapper(cfg, hc)
	if err != nil {
		return nil, err
	}

	return client.New(cfg, client.Options{
		Scheme: scheme,
		Mapper: mapper,
		//Opts: client.WarningHandlerOptions{
		//	SuppressWarnings:   false,
		//	AllowDuplicateLogs: false,
		//},
	})
}

func main() {

	if err := useKubebuilderClient(); err != nil {
		panic(err)
	}
}

func useKubebuilderClient() error {
	fmt.Println("Using kubebuilder client")
	kc, err := NewClient()
	if err != nil {
		return err
	}

	var pglist v1alpha1.FluxCDConfigList
	err = kc.List(context.TODO(), &pglist)
	if err != nil {
		return err
	}
	for _, db := range pglist.Items {
		data, err := yaml.Marshal(db)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		fmt.Println("-------------------------------")

		vals, err := addonfactory.JsonStructToValues(db)
		data2, err := yaml.Marshal(vals)
		if err != nil {
			return err
		}
		fmt.Println(string(data2))
	}

	var list unstructured.UnstructuredList
	list.SetAPIVersion("fluxcd.open-cluster-management.io/v1alpha1")
	list.SetKind("FluxCDConfig")
	err = kc.List(context.TODO(), &list)
	if err != nil {
		return err
	}
	for _, db := range list.Items {
		data, err := yaml.Marshal(db)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		fmt.Println("-------------------------------")
	}

	return nil
}
