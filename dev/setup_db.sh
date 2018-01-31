#!/bin/bash
set -e

echo "Wait for servers to be up"
sleep 10

HOSTPARAMS="--host db --insecure"
SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"

for env in triton triton_test; do
    $SQL -e "CREATE DATABASE IF NOT EXISTS ${env};"
    $SQL -d $env -e "CREATE TABLE IF NOT EXISTS tsg_templates (id SERIAL PRIMARY KEY, name TEXT NOT NULL, package TEXT NOT NULL, image_id TEXT NOT NULL, account_id TEXT NOT NULL, firewall_enabled BOOL DEFAULT false, networks TEXT, metadata TEXT, userdata TEXT, tags TEXT, archived BOOL DEFAULT false);"
    $SQL -d $env -e "CREATE TABLE IF NOT EXISTS tsg_groups (id SERIAL PRIMARY KEY, name TEXT NOT NULL, template_id INT NOT NULL REFERENCES tsg_templates (id), account_id TEXT NOT NULL, capacity INT NOT NULL, datacenter TEXT NOT NULL, health_check_interval INT DEFAULT 300, instance_tags TEXT, archived BOOL DEFAULT false);"
done

/cockroach/cockroach.sh sql $HOSTPARAMS --database=triton < /dev/backup.sql
