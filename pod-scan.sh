#!/bin/bash
## Install and run Trivy and KubeAudit scans against the local K8s cluster
## Assume KubeConfig is correctly configured

# helper function to parse yaml
function parse_yaml {
   local prefix=$2
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}

# install trivy and kubeaudit
echo "Installing Trivy"
sudo rpm -ivh https://github.com/aquasecurity/trivy/releases/download/v0.18.3/trivy_0.18.3_Linux-64bit.rpm
IMAGE=$(parse_yaml config/samples/app_v1_mypodinfo.yaml | grep "spec_image_image" | cut -d '"' -f2)
echo "Installing KubeAudit"
wget https://github.com/Shopify/kubeaudit/releases/download/v0.22.0/kubeaudit_0.22.0_linux_386.tar.gz
tar -xvf kubeaudit_0.22.0_linux_386.tar.gz
rm kubeaudit_0.22.0_linux_386.tar.gz
echo "Finished installing KubeAudit"
echo "Running trivy against $IMAGE"

# run trivy and kubeaudit
trivy image -f table $IMAGE
echo "Completed trivy"
echo "Running KubeAudit"
./kubeaudit all --kubeconfig ~/.kube/config --context minikube
