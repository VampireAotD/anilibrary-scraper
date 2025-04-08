#!/bin/bash

set -eu

sed -i "s|username=\"\"|username=\"${KAFKA_CLIENT_USERS}\"|" /opt/bitnami/kafka/config/client.properties
sed -i "s|password=\"\"|password=\"${KAFKA_CLIENT_PASSWORDS}\"|" /opt/bitnami/kafka/config/client.properties

/opt/bitnami/scripts/kafka/entrypoint.sh "/opt/bitnami/scripts/kafka/run.sh"
