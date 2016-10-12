/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"net/http"
//	"time"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

type Page struct {
    Title string
    Body  []byte
}

func handlerServices(w http.ResponseWriter, r *http.Request) {
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	services, err := clientset.Core().Services("").List(api.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	numberOfServices := fmt.Sprintf("There are %d services in the cluster\n", len(services.Items))
  fmt.Fprintf(w, "<p>Hi there! You requested %s!</p><p>Found %s services.</p>", r.URL.Path[1:], numberOfServices)
}

func main() {
	flag.Parse()

	// p1 := &Page{Title: "Services overview", Body: []byte(numberOfServices)}
	// fmt.Println(string(p1.Body))
	http.HandleFunc("/", handlerServices)
  http.ListenAndServe(":8080", nil)
}
