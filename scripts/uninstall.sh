#!/bin/bash

# Print daemon version
echo "---"
MARKETPLACE_LAST_COMMIT=$(git -C ../marketplace log --oneline -1)
echo "Marketplace last commit:"
echo $MARKETPLACE_LAST_COMMIT

# Load and set environment variables
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

# ===============================================================

# Run uninstall.sh
UNINSTALL_FILE=../marketplace/$CR_APP_NAME/uninstall.sh
if [ -f "$UNINSTALL_FILE" ]
then
    echo "---"
    echo "$UNINSTALL_FILE does exist, running it..."
    chmod +x $UNINSTALL_FILE
    source $UNINSTALL_FILE
    echo "Status of uninstall.sh: $?"
    echo "---"
fi
