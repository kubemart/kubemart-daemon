FROM ubuntu:20.04

ARG HELM_VERSION
ARG KUBECTL_VERSION
RUN apt-get update -y && apt-get install curl -y && \
    # Install helm (https://helm.sh/docs/intro/install)
    curl -L https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz | tar xvz && \
    mv linux-amd64/helm /usr/bin/helm && \
    chmod +x /usr/bin/helm && \
    rm -rf linux-amd64 && \
    # Install kubectl (https://kubernetes.io/docs/tasks/tools/install-kubectl)
    curl -LO https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
    mv kubectl /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl && \
    # Install git (https://git-scm.com/download/linux)
    apt-get install git -y

RUN git clone https://github.com/civo/kubernetes-marketplace.git marketplace
WORKDIR marketplace
