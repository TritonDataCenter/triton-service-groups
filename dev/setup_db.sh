#!/bin/bash
set -e
echo Wait for servers to be up
sleep 10

HOSTPARAMS="--host roach-one --insecure"
SQL="/cockroach/cockroach.sh sql $HOSTPARAMS"

$SQL -e "CREATE DATABASE IF NOT EXISTS triton;"
$SQL -d triton -e "CREATE TABLE IF NOT EXISTS tsg_templates (id SERIAL PRIMARY KEY, name TEXT UNIQUE NOT NULL, package TEXT NOT NULL, image_id TEXT NOT NULL);"