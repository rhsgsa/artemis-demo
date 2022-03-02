FROM openjdk:11

ARG VERSION

RUN \
  curl -Lo /tmp/artemis.zip "https://www.apache.org/dyn/closer.cgi?filename=activemq/activemq-artemis/${VERSION}/apache-artemis-${VERSION}-bin.zip&action=download" \
  && \
  cd /tmp \
  && \
  unzip /tmp/artemis.zip \
  && \
  rm -f /tmp/artemis.zip \
  && \
  mv /tmp/apache-artemis-* /usr/local/artemis
