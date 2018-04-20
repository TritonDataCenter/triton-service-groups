#!/bin/bash

set -e
set -o pipefail

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
    $SQL -e "DROP DATABASE IF EXISTS ${env} CASCADE;"
    $SQL -e "CREATE DATABASE IF NOT EXISTS ${env};"

    cat <<'EOS' | $SQL -d $env
CREATE TABLE IF NOT EXISTS tsg_keys (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" STRING NOT NULL,
    fingerprint STRING NULL,
    material STRING NULL,
    account_id UUID NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    archived BOOL NULL DEFAULT false,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    INDEX name_idx ("name" ASC),
    INDEX id_name_idx (id ASC, "name" ASC),
    INDEX id_account_id_idx (id ASC, account_id ASC),
    INDEX archived_idx (archived ASC),
    FAMILY "primary" (id, "name", fingerprint, material, account_id, created_at, updated_at, archived)
);
EOS

    cat <<'EOS' | $SQL -d $env
CREATE TABLE IF NOT EXISTS tsg_accounts (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    account_name STRING NOT NULL,
    triton_uuid STRING NULL,
    key_id UUID NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    archived BOOL NULL DEFAULT false,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    CONSTRAINT key_id_tsg_keys_id_fk FOREIGN KEY (key_id) REFERENCES tsg_keys (id),
    INDEX key_id_tsg_keys_id_fk_idx (key_id ASC),
    INDEX name_idx (account_name ASC),
    INDEX id_name_idx (id ASC, account_name ASC),
    INDEX archived_idx (archived ASC),
    FAMILY "primary" (id, account_name, triton_uuid, key_id, created_at, updated_at, archived)
);
EOS

    cat <<'EOS' | $SQL -d $env
CREATE TABLE IF NOT EXISTS tsg_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username STRING NOT NULL,
    account_id UUID NOT NULL REFERENCES tsg_accounts (id),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    archived BOOL DEFAULT false
);
EOS

    cat <<'EOS' | $SQL -d $env
CREATE TABLE IF NOT EXISTS tsg_templates (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    template_name STRING NOT NULL,
    account_id UUID NOT NULL,
    package STRING NOT NULL,
    image_id STRING NOT NULL,
    firewall_enabled BOOL NULL DEFAULT false,
    networks STRING NULL,
    userdata STRING NULL,
    metadata STRING NULL,
    tags STRING NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    archived BOOL NULL DEFAULT false,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    CONSTRAINT account_id_tsg_accounts_id_fk FOREIGN KEY (account_id) REFERENCES tsg_accounts (id),
    INDEX account_id_tsg_accounts_id_fk_idx (account_id ASC),
    INDEX name_idx (template_name ASC),
    INDEX archived_idx (archived ASC),
    FAMILY "primary" (id, template_name, account_id, package, image_id, firewall_enabled, networks, userdata, metadata, tags, created_at, archived)
);
EOS

    cat <<'EOS' | $SQL -d $env
CREATE TABLE IF NOT EXISTS tsg_groups (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    "name" STRING NOT NULL,
    template_id UUID NOT NULL,
    account_id UUID NOT NULL,
    capacity INT NOT NULL,
    health_check_interval INT NULL DEFAULT 300:::INT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    archived BOOL NULL DEFAULT false,
    CONSTRAINT "primary" PRIMARY KEY (id ASC),
    CONSTRAINT template_id_tsg_templates_id_fk FOREIGN KEY (template_id) REFERENCES tsg_templates (id),
    CONSTRAINT account_id_tsg_accounts_id_fk FOREIGN KEY (account_id) REFERENCES tsg_accounts (id),
    INDEX template_id_tsg_templates_id_fk_idx (template_id ASC),
    INDEX account_id_tsg_accounts_id_fk_idx (account_id ASC),
    INDEX name_idx ("name" ASC),
    INDEX name_templates_id_idx ("name" ASC, template_id ASC),
    INDEX archived_idx (archived ASC),
    FAMILY "primary" (id, "name", template_id, account_id, capacity, health_check_interval, created_at, updated_at, archived)
);
EOS

    if [ -f /dev/backup.sql ]; then
        /cockroach/cockroach.sh sql $HOSTPARAMS --database=$env < /dev/backup.sql
    fi
done


