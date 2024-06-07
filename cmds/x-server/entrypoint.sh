#!/bin/bash

DISPLAY=:0
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}

/bin/sh -c "until [ -f /tmp/.X0-lock ]; do sleep 0.01; done; exec openbox" &
Xvfb ${DISPLAY} -ac -screen 0 "${SCREEN_RESOLUTION}" -noreset -listen tcp
