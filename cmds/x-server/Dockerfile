FROM ubuntu:22.04

ARG USER_NAME=user
ARG USER_HOME=/home/user
ARG USER_UID=4096
ENV TZ "UTC"

RUN apt update && \
    apt -y install apt-transport-https ca-certificates tzdata locales && \
    DEBIAN_FRONTEND=noninteractive apt upgrade -y && \
    ## install ms fonts
    echo ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula select true | debconf-set-selections && \
    apt install -y xvfb openbox ttf-mscorefonts-installer \
        fontconfig fontconfig-config fonts-dejavu-core \
        fonts-liberation fonts-ubuntu-console fonts-wqy-zenhei \
        fonts-thai-tlwg-ttf fonts-ipafont-mincho fonts-sahadeva \
        fonts-beng fonts-noto-unhinted fonts-noto-color-emoji \
        libfontconfig1 libfontenc1 libfreetype6 \
        libxfont2 libxft2 xfonts-base xfonts-encodings xfonts-utils && \
    ## configure timezone and locales
    ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    mkdir -p /var/lib/locales/supported.d/ && grep UTF-8 /usr/share/i18n/SUPPORTED > /var/lib/locales/supported.d/all && \
    locale-gen && update-locale && \
    ## clean font caches
    fc-cache -f -v && \
    apt clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*

RUN adduser --system \
    --home $USER_HOME \
    --uid $USER_UID \
    --ingroup root \
    --disabled-password \
    --shell /bin/bash \
    $USER_NAME

USER $USER_NAME

COPY menu.xml /etc/xdg/openbox/menu.xml
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]