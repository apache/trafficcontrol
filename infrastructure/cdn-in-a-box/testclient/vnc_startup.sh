#!/bin/bash

VNC_DEPTH=${VNC_DEPTH:-32}
VNC_RESOLUTION=${VNC_DEPTH:-1440x900}
VNC_WM=${VNC_WM:-fluxbox}

vncconfig -iconic &
firefox &
xterm -bg black -fg white +sb &

case $VNC_WM in
  "xfce")
    startxfce4 &
    ;;
  "fluxbox")
    startfluxbox &
    ;;
  *)
    echo "No such Window Manager"
  ;;
esac
