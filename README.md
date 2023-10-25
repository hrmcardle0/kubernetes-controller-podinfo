# PodInfo Kubernetes Controller 

This project demonstrates how to create a simple k8s controller based on [KubeBuilder](https://github.com/kubernetes-sigs/kubebuilder). This controller was already pre-built based on the documentation. 
The docs are simple to follow and include a few steps including the init'ing a new project, creating a custom resource, installing that custom resource, editing the type definition, and generating a controller.

The controller's main procedure is the Go reconcile function that continously runs and responds to CRUD events based on your custom resource defition. In this project, the custom resource we are creating
is a pod based on the [PodInfo](https://github.com/stefanprodan/podinfo) project. This is a simple microservice application deployed as a pod with a UI and API.

# Installation

Pre-reqs:
- Docker (authenticated)
- Minikube or some other cluster
- Kubectl
- 'chmod +x ./controller-setup' is executed

This project assumes you can successfuly deploy pods to your k8s cluster and your kubeconfig is correctly set, either by default path ~/.kube/config or the KUBECONFIG env variable

Running the script:

`./controller-setup.sh -i hrmcardle0/mypodinfo:latest`

This script assumes you have been logged into docker to perform the pull. It is normal to see a reosurce-denied error when the push is attempted via the Makefile, as the script hasn't authenticated. 
This will naturally just not push the image, but the script continues. 

The script then applies a custom cluster role & binding giving the system:serviceaccount:mypodinfo-system:mypodinfo-controller-manager access to perform it's duties from within the cluster.
A NodePort service is then started as a front-end for the podinfo app.

The controller is now installed. 