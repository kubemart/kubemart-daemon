#!/bin/bash

echo "---"

# load and set environment variables
source .env

# Read JSON file
FILE_CONTENT=$(cat installation-info.json)
CR_APP_NAME=$(echo $FILE_CONTENT | jq -r ".cr_name")
CR_APP_NAMESPACE=$(echo $FILE_CONTENT | jq -r ".cr_namespace")
echo "CR APP NAME      :" $CR_APP_NAME
echo "CR APP NAMESPACE :" $CR_APP_NAMESPACE

# Capture installation status so we can proceed 
# with post-install steps (update the CR status)
lastStatus=1

# Run pre-install.sh
PRE_INSTALL_FILE=../marketplace/$CR_APP_NAME/pre_install.sh
if [ -f "$PRE_INSTALL_FILE" ]
then
    echo "---"
    echo "$PRE_INSTALL_FILE does exist, running it..."
    chmod +x $PRE_INSTALL_FILE
    source $PRE_INSTALL_FILE
    lastStatus=$?
fi

# Run kubectl againts app.yaml
APP_YAML_FILE=../marketplace/$CR_APP_NAME/app.yaml
if [ -f "$APP_YAML_FILE" ]
then
    echo "---"
    echo "$APP_YAML_FILE does exist, running it (using kubectl)..."
    cat app.yaml | envsubst | kubectl apply -f -
    lastStatus=$?
fi

# Run install.sh
INSTALL_FILE=../marketplace/$CR_APP_NAME/install.sh
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

    # TODO
    # Update installed version to App CR
else
    echo "---"
    echo "Something went wrong with kubectl command..."
    exit 1
fi
