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
	"html/template"

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

const serviceTmpl = `
<div>
	<h2><a href='{{.URL}}'>Service '{{.serviceName}}'</a></h2>
	<p><ul>
		<li>URL: {{.URL}}
		<li>Port: {{.port}}
		<li><a href='{{.k8sAPI}}'>Kubernetes API object: '{{.k8sAPI}}'</a>
	</ul></p>
</div>
`

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
	services_list, err := clientset.Core().Services("").List(api.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	nodes_list, err := clientset.Core().Nodes().List(api.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	nodeAddress := nodes_list.Items[0].ObjectMeta.Name
  fmt.Fprintf(w, "<h1>Hi there!<br />Found %d services.</h1>\n", len(services_list.Items))
	for _,service := range services_list.Items {
		if service.Spec.Ports[0].NodePort > 0 {
			nodePort := service.Spec.Ports[0].NodePort
			serviceURL := fmt.Sprintf("http://%s:%d/", nodeAddress, nodePort)
			serviceName := service.ObjectMeta.Name
			k8sAPI := fmt.Sprintf("http://%s:8080%s", nodeAddress, service.ObjectMeta.SelfLink)
			serviceData := map[string]interface{}{
		    "URL":         serviceURL,
		    "port":        nodePort,
		    "serviceName": serviceName,
		    "k8sAPI":      k8sAPI,
			}
			t := template.Must(template.New("service").Parse(serviceTmpl))
			if err := t.Execute(w, serviceData); err != nil {
		    panic(err)
			}
		}
	}
}

func main() {
	flag.Parse()

	// p1 := &Page{Title: "Services overview", Body: []byte(numberOfServices)}
	// fmt.Println(string(p1.Body))
	http.HandleFunc("/", handlerServices)
  http.ListenAndServe(":8080", nil)
}
