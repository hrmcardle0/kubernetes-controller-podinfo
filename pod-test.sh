#!/bin/bash
## Deply and run various tests against the cluster using the CRD specified
## Assume KubeConfig is correctly configured

function parse_yaml() {
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


# arg parse
while [[ $# -gt 0 ]]; do
    case $1 in 
        -c|--crd)
            FILE="$2"
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

## tests
CONTROLLER_INSTALLED=failed
POD_SUCCESSFUL=failed
UPDATE_SUCCESSFUL=failed
CURL_SUCCESSFUL=failed

# define file name
DEF=$(echo $FILE | cut -d '/' -f3)
echo "Starting testing script based on CRD located at: $FILE"

## ensure controller is installed
kubectl get pods -n mypodinfo-system  | grep controller > /dev/null
if [ $? -eq 0 ]; then
    CONTROLLER_INSTALLED=passed
fi

## start pod and ensure successful
kubectl apply -f $FILE > /dev/null
if [ $? -eq 0 ]; then
    POD_SUCCESSFUL=passed
fi

# wait then curl endpoint
## get pod
POD=$(kubectl get pod -l app=podinfo -o jsonpath="{.items[0].metadata.name}")

## loop until pod is running
POD_STATUS=$(kubectl get pod $POD -o jsonpath="{.status.phase}")
while [ $POD_STATUS != "Running" ]
do
    echo "Pod is not running. Waiting 5 seconds..."
    sleep 5
    POD_STATUS=$(kubectl get pod $POD -o jsonpath="{.status.phase}")
done
echo "Pod is running, attempting curl of service"
curl -sd 'anon' 192.168.49.2:30163/token | grep token
if [ $? -eq 0 ]; then
    CURL_SUCCESSFUL=passed
fi 

# update pod and ensure successful
sed -i 's/replicaCount: 1/replicaCount: 2/' $FILE
kubectl apply -f $FILE >/dev/null
sleep 2
if [ $(kubectl get rc podinfo -o jsonpath="{.status.replicas}") -eq 2 ]; then
    UPDATE_SUCCESSFUL=passed
fi

# reverting update
sed -i 's/replicaCount: 2/replicaCount: 1/' $FILE

echo "1) Controller Test Result: $CONTROLLER_INSTALLED"
echo "2) Pod Test Result: $POD_SUCCESSFUL"
echo "3) Curl Test Result: $CURL_SUCCESSFUL"
echo "4) Update Test Result: $UPDATE_SUCCESSFUL"

