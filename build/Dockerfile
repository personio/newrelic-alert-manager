FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/newrelic-alert-manager \
    USER_UID=1001 \
    USER_NAME=newrelic-alert-manager

# install operator binary
COPY build/_output/bin/newrelic-alert-manager ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
