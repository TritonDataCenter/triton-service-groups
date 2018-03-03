//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package templates_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx"
)

func FindByName(name string, db *pgx.ConnPool, accountId string) (*InstanceTemplate, bool) {
	var (
		template     InstanceTemplate
		metaDataJson string
		tagsJson     string
		networksList string
	)

	sqlStatement := `SELECT id, template_name, package, image_id, instance_name_prefix, account_id, firewall_enabled, networks, COALESCE(metadata,''), userdata, COALESCE(tags,'')
FROM triton.tsg_templates
WHERE template_name = $1 and account_id = $2
AND archived = false;`

	err := db.QueryRowEx(context.TODO(), sqlStatement, nil, name, accountId).
		Scan(&template.ID,
			&template.TemplateName,
			&template.Package,
			&template.ImageId,
			&template.InstanceNamePrefix,
			&template.AccountId,
			&template.FirewallEnabled,
			&networksList,
			&metaDataJson,
			&template.UserData,
			&tagsJson)
	switch err {
	case nil:
		metaData, err := convertFromJson(metaDataJson)
		if err != nil {
			panic(err)
		}
		template.MetaData = metaData

		tags, err := convertFromJson(tagsJson)
		if err != nil {
			panic(err)
		}
		template.Tags = tags

		template.Networks = strings.Split(networksList, ",")

		return &template, true
	case pgx.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil, false
	default:
		panic(err)
	}
}

func FindAll(db *pgx.ConnPool, accountId string) ([]*InstanceTemplate, error) {
	var templates []*InstanceTemplate

	sqlStatement := `SELECT id, template_name, package, image_id, account_id, firewall_enabled, instance_name_prefix, networks, COALESCE(metadata,''), userdata, COALESCE(tags, '')
FROM triton.tsg_templates
WHERE account_id = $1
AND archived = false;`

	var metaDataJson string
	var tagsJson string
	var networksList string

	rows, err := db.QueryEx(context.TODO(), sqlStatement, nil, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var template InstanceTemplate
		err := rows.Scan(&template.ID,
			&template.TemplateName,
			&template.Package,
			&template.ImageId,
			&template.AccountId,
			&template.FirewallEnabled,
			&template.InstanceNamePrefix,
			&networksList,
			&metaDataJson,
			&template.UserData,
			&tagsJson)
		if err != nil {
			return nil, err
		}

		metaData, err := convertFromJson(metaDataJson)
		if err != nil {
			panic(err)
		}
		template.MetaData = metaData

		tags, err := convertFromJson(tagsJson)
		if err != nil {
			panic(err)
		}
		template.Tags = tags

		template.Networks = strings.Split(networksList, ",")

		templates = append(templates, &template)
	}

	return templates, nil
}

func (t *InstanceTemplate) Save(db *pgx.ConnPool, accountId string) {
	sqlStatement := `
INSERT INTO triton.tsg_templates (template_name, package, image_id, account_id, firewall_enabled, instance_name_prefix, networks, metadata, userdata, tags)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

	metaDataJson, err := convertToJson(t.MetaData)
	if err != nil {
		log.Fatal(err)
	}

	tagsJson, err := convertToJson(t.Tags)
	if err != nil {
		log.Fatal(err)
	}

	networksList := strings.Join(t.Networks, ",")

	_, err = db.ExecEx(context.TODO(), sqlStatement, nil,
		t.TemplateName, t.Package, t.ImageId,
		accountId, t.FirewallEnabled, t.InstanceNamePrefix, networksList, metaDataJson,
		t.UserData, tagsJson)
	if err != nil {
		panic(err)
	}
}

func (t *InstanceTemplate) Archive(db *pgx.ConnPool, accountId string) {
	sqlStatement := `UPDATE triton.tsg_templates
SET archived = true
WHERE template_name = $1 and account_id = $2`

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil, t.TemplateName, accountId)
	if err != nil {
		panic(err)
	}
}

func convertToJson(data map[string]string) (string, error) {
	if data == nil {
		return "", nil
	}

	log.Printf("Found data")
	json, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(json), nil
}

func convertFromJson(data string) (map[string]string, error) {
	if data == "" {
		return nil, nil
	}

	var result map[string]string
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
