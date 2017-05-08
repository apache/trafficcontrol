for f in /opt/traffic_ops/install/data/profiles/*.sql ; do psql -U traffic_ops -h localhost traffic_ops -f $f ; done
