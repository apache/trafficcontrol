#!/bin/bash

################################################################################
# Wait on SSL certificate generation
until [ -f "$X509_CA_DONE_FILE" ] 
do
  echo "Waiting on Shared SSL certificate generation"
  sleep 3
done

# Source the CIAB-CA shared SSL environment
source $X509_CA_ENV_FILE

# Trust the CIAB-CA at the System level
cp $X509_CA_CERT_FILE /etc/pki/ca-trust/source/anchors
update-ca-trust extract
################################################################################

VNC_DEPTH=${VNC_DEPTH:-32}
VNC_RESOLUTION=${VNC_RESOLUTION:-1440x900}

su -c "vncserver :0 -depth $VNC_DEPTH -geometry $VNC_RESOLUTION" - "$VNC_USER" && tail -F /home/dev/.vnc/vnc:0.log
