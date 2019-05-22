FROM registry.access.redhat.com/ubi7-dev-preview/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/orca-operator \
    USER_UID=1001 \
    USER_NAME=orca-operator

ARG version
ARG release

LABEL name=orca-operator
LABEL vendor=Tufin
LABEL version=$version
LABEL release=$release
LABEL summary="Tufin Orca is a security solution that enables automation of network security policies for applications running on Kubernetes platforms."
LABEL description="Orca Operator deploys the Orca agent into a Kubernetes cluster."

# install operator binary
COPY build/_output/bin/orca-operator ${OPERATOR}
COPY licenses/ /licenses/

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
