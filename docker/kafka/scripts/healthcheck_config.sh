#!/bin/bash

if [ -n "$KAFKA_CLIENT_USERS" ] && [ -n "$KAFKA_CLIENT_PASSWORDS" ]; then
  sed -i "s|username=\"\"|username=\"${KAFKA_CLIENT_USERS}\"|" /opt/bitnami/kafka/config/client.properties
  sed -i "s|password=\"\"|password=\"${KAFKA_CLIENT_PASSWORDS}\"|" /opt/bitnami/kafka/config/client.properties
fi

/opt/bitnami/scripts/kafka/entrypoint.sh "/run.sh"
