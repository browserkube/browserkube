#!/bin/sh

PORT=${PORT:-"5900"}
DISPLAY=${DISPLAY:-":0"}
DISPLAY_ARG=${DISPLAY_ARG:-"WAIT:127.0.0.1:0"}
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
