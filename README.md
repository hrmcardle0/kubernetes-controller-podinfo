# PodInfo Kubernetes Controller 

This project demonstrates how to create a simple k8s controller based on [KubeBuilder](https://github.com/kubernetes-sigs/kubebuilder). This controller was already pre-built based on the documentation. 
The docs are simple to follow and include a few steps including the init'ing a new project, creating a custom resource, installing that custom resource, editing the type definition, and generating a controller.

The controller's main procedure is the Go reconcile function that continously runs and responds to CRUD events based on your custom resource defition. In this project, the custom resource we are creating
is a replicationcontroller based on the [PodInfo](https://github.com/stefanprodan/podinfo) project. This is a simple microservice application deployed as a pod with a UI and API.

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

## Usage

The sample CRD is contained within the file: config/samples/app_v1_mypodinfo.yaml

It can then be edited, updated and deleted, which cooresponds to the podinfo resource itself.

Tests can be run via the following command, which exexcutes the following tests:

- Ensure the controller is registered
- Ensure creating an instance of the CRD correctly creates a pod of type podinfo
- Ensure the podinfo application is correctly bootstrapped by curling the token endpoint
- Ensure the custom resource can be updated and corresponding changes will update the underlying infrastructure appropriately
- Ensure the custom resource is correctly deleted with all corresponding infrastructure deleted

`./pod-test.sh -c config/samples/app_v1_mypodinfo.yaml`


This script uses trivy and kubeaudit to run security checks against the podinfo image as well as the cluster itself:

`./pod-scan`

## Cleanup

To clean up your environment, remove the custom resource:

`kubectl delete mypodinfo --all`

Undeploy the controller:

`make undeploy`

## Troubleshooting

When updating behind 1 replicas, 1 pod is running while the rest are stuck in pending

- This implementation was implemented on minikube and thus binds staticly to the host port. Another type of service or multiple nodes will need to be created for multiple pods to be running

When creating the custom resource with kubectl apply, I don't see any pods created

- The controller is not running, ensure the output of the controller-setup script correctly shows the image being pulled and the controller is installed. 