

ARG BASE_IMAGE=rockylinux \
    RHEL_VERSION=8
FROM ${BASE_IMAGE}:${RHEL_VERSION}
ARG RHEL_VERSION=8

RUN if [[ "${RHEL_VERSION%%.*}" -eq 7 ]]; then \
        yum -y install dnf || exit 1; \
    fi


#ADD traffic_portal_v2/analytics/kubernetes /
#RUN mv /kubernetes /etc/yum.repos.d/kubernetes.repo
RUN dnf install -y epel-release && \
    dnf install -y openssl bind-utils net-tools gettext findutils curl nfs-utils \
    php php-curl php-gd php-json mysql mysql-server php-xml php-mbstring php-pdo php-mysqli \
    ansible git unzip wget nginx

ADD traffic_portal_v2/analytics/run.sh \
  traffic_portal_v2/analytics/matomo.conf \
  traffic_portal_v2/analytics/mysql.sql \
  traffic_portal_v2/analytics/config.ini.php \
  traffic_portal_v2/analytics/matomo.sql \
  traffic_ops/to-access.sh \
  enroller/server_template.json /

COPY dns/set-dns.sh \
     dns/insert-self-into-dns.sh \
     /usr/local/sbin/

EXPOSE 443

CMD /run.sh
