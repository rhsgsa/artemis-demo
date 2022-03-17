FROM openjdk:11

ARG VERSION

RUN \
  apt update \
  && \
  apt install -y patch dnsutils netcat inetutils-ping \
  && \
  apt clean \
  && \
  curl -Lo /tmp/artemis.zip "https://www.apache.org/dyn/closer.cgi?filename=activemq/activemq-artemis/${VERSION}/apache-artemis-${VERSION}-bin.zip&action=download" \
  && \
  cd /tmp \
  && \
  unzip /tmp/artemis.zip \
  && \
  rm -f /tmp/artemis.zip \
  && \
  mv /tmp/apache-artemis-* /usr/local/artemis \
  && \
  useradd -m -u 1001 -g 0 demo \
  && \
  mkdir -p /storage/instances/node1 /storage/instances/node2 /storage/shared-store \
  && \
  chown -R 1001:0 /storage \
  && \
  echo 'set -o vi' >> /home/demo/.bashrc

USER 1001
