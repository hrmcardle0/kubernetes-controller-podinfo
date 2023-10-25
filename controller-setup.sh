#!/bin/bash

# arg parse
while [[ $# -gt 0 ]]; do
    case $1 in 
        -i|--image)
            echo "asdfasd"
            IMAGE="$2"
            shift
            shift
            ;;
        -*)
            echo "Unknown option $1"
            shift
            shift
            ;;
    esac
done
echo "Bootstrapping environment for MyPodinfo Controller"
echo "(make) Generating Manfiest"
make manifests
echo "(make) Installing CRD"
make install
echo "(make) Build docker image...$IMAGE"
make docker-build docker-push IMG=$IMAGE
echo "Applying custom role"
kubectl apply -f custom/cluster-role.yaml
kubectl apply -f custom/cluster-role-bindings.yaml
echo "Applying service"
kubectl apply -f custom/service.yaml
echo "(make) Deploying docker image...$IMAGE"
make deploy IMG=$IMAGE