# ====================================================
# Go app
# ====================================================
# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o main main.go

# ====================================================
# Others
# ====================================================
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
    apt-get install git -y && \
    # Install jq (https://stedolan.github.io/jq/download/)
    apt-get install jq -y

WORKDIR /
ADD scripts /scripts
RUN chmod +x /scripts/run.sh
RUN git clone https://github.com/civo/kubernetes-marketplace.git marketplace
COPY --from=builder /workspace/main .
