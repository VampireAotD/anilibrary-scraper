ARG KAFKA_VERSION=latest

FROM bitnami/kafka:$KAFKA_VERSION

COPY config/client.properties /opt/bitnami/kafka/config/client.properties

COPY --chown=1001:1001 scripts/client_entrypoint.sh /opt/bitnami/scripts/kafka/client_entrypoint.sh

RUN chmod +x /opt/bitnami/scripts/kafka/client_entrypoint.sh

ENTRYPOINT ["/opt/bitnami/scripts/kafka/client_entrypoint.sh"]
