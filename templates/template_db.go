//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package templates_v1

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/jackc/pgx"
)

func FindTemplateByID(db *pgx.ConnPool, key int64, accountId string) (*InstanceTemplate, bool) {
	var template InstanceTemplate

	sqlStatement := `SELECT id, template_name, package, image_id, instance_name_prefix, account_id, firewall_enabled, networks, COALESCE(metadata,''), userdata, COALESCE(tags,'')  
FROM triton.tsg_templates 
WHERE id = $1 and account_id = $2
AND archived = false;`

	var metaDataJson string
	var tagsJson string
	var networksList string

	err := db.QueryRowEx(context.TODO(), sqlStatement, nil, key, accountId).
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
		return nil, false
	default:
		panic(err)
	}
}

func FindTemplates(db *pgx.ConnPool, accountId string) ([]*InstanceTemplate, error) {
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

func SaveTemplate(db *pgx.ConnPool, accountId string, template *InstanceTemplate) {
	sqlStatement := `
INSERT INTO triton.tsg_templates (template_name, package, image_id, account_id, firewall_enabled, instance_name_prefix, networks, metadata, userdata, tags) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

	metaDataJson, err := convertToJson(template.MetaData)
	if err != nil {
		log.Fatal(err)
	}

	tagsJson, err := convertToJson(template.Tags)
	if err != nil {
		log.Fatal(err)
	}

	networksList := strings.Join(template.Networks, ",")

	_, err = db.ExecEx(context.TODO(), sqlStatement, nil,
		template.TemplateName, template.Package, template.ImageId,
		accountId, template.FirewallEnabled, template.InstanceNamePrefix, networksList, metaDataJson,
		template.UserData, tagsJson)
	if err != nil {
		panic(err)
	}
}

func RemoveTemplate(db *pgx.ConnPool, identifier int64, accountId string) {
	sqlStatement := `UPDATE triton.tsg_templates 
SET archived = true 
WHERE id = $1 and account_id = $2`

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil, identifier, accountId)
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
