FROM alpine:3.19

ARG USER_NAME=user
ARG USER_HOME=/home/user
ARG USER_UID=4096

RUN \
    apk add --update --no-cache x11vnc

RUN adduser --system \
    --home $USER_HOME \
    --uid $USER_UID \
    --ingroup root \
    --disabled-password \
    --shell /bin/bash \
    $USER_NAME

USER $USER_NAME

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]