#!/bin/sh

PORT=${PORT:-"5900"}
DISPLAY=${DISPLAY:-":0"}
DISPLAY_ARG=${DISPLAY_ARG:-"WAIT:50:${DISPLAY}"}
VNCPASS=${VNCPASS:-"browserkube"}
x11vnc \
  -xkb \
  -xrandr \
  -passwd "${VNCPASS}" \
  -noxrecord \
  -forever \
  -display "${DISPLAY_ARG}" \
  -shared \
  -rfbport $PORT
