FROM golang:1.24

ENV GO111MODULE=on
ENV XDG_CACHE_HOME=/tmp/.cache
ENV GOPATH=/tmp/go
ENV GOBIN=/tmp/gobin

ARG KUBERNETES_BRANCH="master"
ARG REPO="${GOPATH}/src/github.com/clbiggs/cyberark-to-k8s"

RUN DEBIAN_FRONTEND=noninteractive apt-get update \
  && DEBIAN_FRONTEND=noninteractive apt install -y bash git openssh-client \
  && git clone --depth 1 --branch ${KUBERNETES_BRANCH} \
  https://github.com/kubernetes/code-generator.git ${GOPATH}/src/k8s.io/code-generator

RUN cd ${GOPATH}/src/k8s.io/code-generator \
  && go mod download \
  && mkdir -p ${GOBIN} \
  && mkdir -p ${REPO} \
  && find ${GOPATH} -type d -exec chmod 755 {} \; \
  && find ${GOPATH} -type f -exec chmod 666 {} \; \
  && chmod +x ${GOPATH}/src/k8s.io/code-generator/kube_codegen.sh

WORKDIR ${REPO}

VOLUME ${REPO}
