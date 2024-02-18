ARG KAFKA_VERSION

FROM bitnami/kafka:$KAFKA_VERSION

COPY config/client.properties /opt/bitnami/kafka/config/client.properties

COPY scripts/healthcheck_config.sh /opt/bitnami/scripts/kafka/healthcheck_config.sh

ENTRYPOINT ["/opt/bitnami/scripts/kafka/healthcheck_config.sh"]
