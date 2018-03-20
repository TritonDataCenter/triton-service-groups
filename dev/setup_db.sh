#!/bin/bash
set -e

if [ -x "/cockroach/cockroach.sh" ]; then
    echo "Wait for servers to be up"
    sleep 10

    HOSTPARAMS="--host db --insecure"
    SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"
else
    HOSTPARAMS="--host localhost --insecure --certs-dir ./dev/vagrant/certs"
    SQL="cockroach sql $HOSTPARAMS"
fi

for env in triton triton_test; do
    $SQL -e "CREATE DATABASE IF NOT EXISTS ${env};"

    $SQL -d $env -e "CREATE TABLE IF NOT EXISTS tsg_accounts (
id SERIAL PRIMARY KEY,
account_name STRING NOT NULL,
triton_uuid STRING,
key_name STRING,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
archived BOOL DEFAULT false);"

    $SQL -d $env -e "CREATE TABLE IF NOT EXISTS tsg_users (
id SERIAL PRIMARY KEY,
username STRING NOT NULL,
account_id INT NOT NULL REFERENCES tsg_accounts (id),
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
archived BOOL DEFAULT false);"

    $SQL -d $env -e "CREATE TABLE IF NOT EXISTS tsg_templates (
id SERIAL PRIMARY KEY,
template_name STRING NOT NULL,
account_id INT NOT NULL REFERENCES tsg_accounts (id),
package STRING NOT NULL,
image_id STRING NOT NULL,
instance_name_prefix STRING,
firewall_enabled BOOL DEFAULT false,
networks TEXT,
userdata TEXT,
metadata TEXT,
tags TEXT,
archived BOOL DEFAULT false);"

    $SQL -d $env -e "CREATE TABLE IF NOT EXISTS tsg_groups (
id SERIAL PRIMARY KEY,
name STRING NOT NULL,
template_id INT NOT NULL REFERENCES tsg_templates (id),
account_id INT NOT NULL REFERENCES tsg_accounts (id),
capacity INT NOT NULL,
health_check_interval INT DEFAULT 300,
archived BOOL DEFAULT false);"

    if [ -f /dev/backup.sql ]; then
        /cockroach/cockroach.sh sql $HOSTPARAMS --database=$env < /dev/backup.sql
    fi
done


