#!/bin/bash

# load and set environment variables
echo "---"
ENV_FILE=.env
if [ -f "$ENV_FILE" ]
then
    cat $ENV_FILE
    source $ENV_FILE
else
    echo "This app does not have .env file (meaning it does not have configurations)"
fi
echo "---"

# Read JSON file
echo "---"
FILE_CONTENT=$(cat installation-info.json)
CR_APP_NAME=$(echo $FILE_CONTENT | jq -r ".cr_name")
CR_APP_NAMESPACE=$(echo $FILE_CONTENT | jq -r ".cr_namespace")
echo "CR APP NAME      :" $CR_APP_NAME
echo "CR APP NAMESPACE :" $CR_APP_NAMESPACE
echo "---"

# Run pre-install.sh
PRE_INSTALL_FILE=../marketplace/$CR_APP_NAME/pre_install.sh
if [ -f "$PRE_INSTALL_FILE" ]
then
    echo "---"
    echo "$PRE_INSTALL_FILE does exist, running it..."
    chmod +x $PRE_INSTALL_FILE
    source $PRE_INSTALL_FILE
    echo "Status of pre_install.sh: $?"
    echo "---"
fi

# Run kubectl againts app.yaml
APP_YAML_FILE=../marketplace/$CR_APP_NAME/app.yaml
if [ -f "$APP_YAML_FILE" ]
then
    echo "---"
    echo "$APP_YAML_FILE does exist, running it (using kubectl)..."
    cat $APP_YAML_FILE | envsubst | kubectl apply -f -
    echo "Status of app.yaml: $?"
    echo "---"
fi

# Run install.sh
INSTALL_FILE=../marketplace/$CR_APP_NAME/install.sh
if [ -f "$INSTALL_FILE" ]
then
    echo "---"
    echo "$INSTALL_FILE does exist, running it..."
    chmod +x $INSTALL_FILE
    source $INSTALL_FILE
    echo "Status of install.sh: $?"
    echo "---"
fi
