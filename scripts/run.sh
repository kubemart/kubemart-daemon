#!/bin/bash
echo "---"

# Read JSON file
FILE_CONTENT=$(cat installation-info.json)
APP_NAME=$(echo $FILE_CONTENT | jq -r ".AppName") # app name from marketplace i.e. wordpress
CRD_NAME=$(echo $FILE_CONTENT | jq -r ".CrdName") # metadata.name of App CRD i.e. my-wordpress
NAMESPACE=$(echo $FILE_CONTENT | jq -r ".Namespace")
HELM_RELEASE_NAME=$(echo $FILE_CONTENT | jq -r ".HelmReleaseName")
HELM_METADATA_NAME=$(echo $FILE_CONTENT | jq -r ".HelmMetadataName")
echo "APP_NAME           :" $APP_NAME
echo "CRD_NAME           :" $CRD_NAME
echo "NAMESPACE          :" $NAMESPACE
echo "HELM_RELEASE_NAME  :" $HELM_RELEASE_NAME
echo "HELM_METADATA_NAME :" $HELM_METADATA_NAME

# Get the most recent Git commit hash of the app.yaml file from marketplace folder
TRACKING_FILE=""
COMMIT_HASH=""
COMMIT_HASH=$(git --git-dir ../marketplace/.git log -n 1 --pretty=format:%H -- $APP_NAME/app.yaml)
TRACKING_FILE="app.yaml"
if [ -z "$COMMIT_HASH" ]
then
    COMMIT_HASH=$(git --git-dir ../marketplace/.git log -n 1 --pretty=format:%H -- $APP_NAME/install.sh)
    TRACKING_FILE="install.sh"
fi
echo "COMMIT_HASH        :" $COMMIT_HASH
echo "TRACKING_FILE      :" $TRACKING_FILE

# Capture installation status so we can proceed with post-install steps
lastStatus=1

# Run pre-install.sh
PRE_INSTALL_FILE=../marketplace/$APP_NAME/pre_install.sh
if [ -f "$PRE_INSTALL_FILE" ]
then
    echo "---"
    echo "$PRE_INSTALL_FILE does exist, running it..."
    chmod +x $PRE_INSTALL_FILE
    source $PRE_INSTALL_FILE
    lastStatus=$?
fi

# Run kubectl againts YAML files
APP_YAML_FILE=../marketplace/$APP_NAME/app.yaml
if [ -f "$APP_YAML_FILE" ]
then
    echo "---"
    echo "Applying YAML files using kubectl..."
    kubectl apply $(ls *.yaml | awk ' { print " -f " $1 } ')
    lastStatus=$?
fi

# Run install.sh
INSTALL_FILE=../marketplace/$APP_NAME/install.sh
if [ -f "$INSTALL_FILE" ]
then
    echo "---"
    echo "$INSTALL_FILE does exist, running it..."
    chmod +x $INSTALL_FILE
    source $INSTALL_FILE
    lastStatus=$?
fi

if [ "$lastStatus" -eq 0 ]; then
    echo "---"
    echo "Updating CRD..."
    TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)

    echo "---"
    echo "Update Git commit hash..."
    curl --silent --show-error --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt --location --request PATCH \
     "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_PORT_443_TCP_PORT/apis/app.bizaar.civo.com/v1alpha1/namespaces/$NAMESPACE/apps/$CRD_NAME/" \
        --header "Content-Type: application/json-patch+json" \
        --header "Authorization: Bearer $TOKEN" \
        --data-raw '[{"op": "replace","path": "/spec/githash","value": "'"$COMMIT_HASH"'"}]'

    echo "---"
    echo "Update tracking file..."
    curl --silent --show-error --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt --location --request PATCH \
     "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_PORT_443_TCP_PORT/apis/app.bizaar.civo.com/v1alpha1/namespaces/$NAMESPACE/apps/$CRD_NAME/" \
        --header "Content-Type: application/json-patch+json" \
        --header "Authorization: Bearer $TOKEN" \
        --data-raw '[{"op": "replace","path": "/spec/gitfile","value": "'"$TRACKING_FILE"'"}]'
    
    if [[ $HELM_RELEASE_NAME != "null" ]]
    then
        echo "---"
        echo "Update Helm release name..."
        curl --silent --show-error --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt --location --request PATCH \
        "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_PORT_443_TCP_PORT/apis/app.bizaar.civo.com/v1alpha1/namespaces/$NAMESPACE/apps/$CRD_NAME/" \
            --header "Content-Type: application/json-patch+json" \
            --header "Authorization: Bearer $TOKEN" \
            --data-raw '[{"op": "replace","path": "/spec/helmreleasename","value": "'"$HELM_RELEASE_NAME"'"}]'
    fi
    
    if [[ $HELM_METADATA_NAME != "null" ]]
    then
        echo "---"
        echo "Update Helm metadata.name..."
        curl --silent --show-error --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt --location --request PATCH \
        "https://$KUBERNETES_SERVICE_HOST:$KUBERNETES_PORT_443_TCP_PORT/apis/app.bizaar.civo.com/v1alpha1/namespaces/$NAMESPACE/apps/$CRD_NAME/" \
            --header "Content-Type: application/json-patch+json" \
            --header "Authorization: Bearer $TOKEN" \
            --data-raw '[{"op": "replace","path": "/spec/helmmetadataname","value": "'"$HELM_METADATA_NAME"'"}]'
    fi
else
    echo "---"
    echo "Something went wrong with kubectl command..."
    exit 1
fi
