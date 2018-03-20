//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package templates_v1

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/jackc/pgx"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/rs/zerolog/log"
)

func FindTemplateByName(ctx context.Context, key string, accountID int) (*InstanceTemplate, bool) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return nil, false
	}

	var template InstanceTemplate

	sqlStatement := `SELECT id, template_name, package, image_id, instance_name_prefix, account_id, firewall_enabled, networks, COALESCE(metadata,''), userdata, COALESCE(tags,'')
FROM triton.tsg_templates
WHERE template_name = $1 and account_id = $2
AND archived = false;`

	var metaDataJson string
	var tagsJson string
	var networksList string

	err := db.QueryRowEx(ctx, sqlStatement, nil, key, accountID).
		Scan(&template.ID,
			&template.TemplateName,
			&template.Package,
			&template.ImageID,
			&template.InstanceNamePrefix,
			&template.AccountID,
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

func FindTemplateByID(ctx context.Context, key int64, accountID int) (*InstanceTemplate, bool) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return nil, false
	}

	var template InstanceTemplate

	sqlStatement := `SELECT id, template_name, package, image_id, instance_name_prefix, account_id, firewall_enabled, networks, COALESCE(metadata,''), userdata, COALESCE(tags,'')
FROM triton.tsg_templates
WHERE id = $1 and account_id = $2
AND archived = false;`

	var metaDataJson string
	var tagsJson string
	var networksList string

	err := db.QueryRowEx(ctx, sqlStatement, nil, key, accountID).
		Scan(&template.ID,
			&template.TemplateName,
			&template.Package,
			&template.ImageID,
			&template.InstanceNamePrefix,
			&template.AccountID,
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

func FindTemplates(ctx context.Context, accountID int) ([]*InstanceTemplate, error) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return nil, handlers.ErrNoConnPool
	}

	var templates []*InstanceTemplate

	sqlStatement := `SELECT id, template_name, package, image_id, account_id, firewall_enabled, instance_name_prefix, networks, COALESCE(metadata,''), userdata, COALESCE(tags, '')
FROM triton.tsg_templates
WHERE account_id = $1
AND archived = false;`

	var metaDataJson string
	var tagsJson string
	var networksList string

	rows, err := db.QueryEx(ctx, sqlStatement, nil, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var template InstanceTemplate
		err := rows.Scan(&template.ID,
			&template.TemplateName,
			&template.Package,
			&template.ImageID,
			&template.AccountID,
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

func SaveTemplate(ctx context.Context, accountID int, template *InstanceTemplate) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return
	}

	sqlStatement := `
INSERT INTO triton.tsg_templates (template_name, package, image_id, account_id, firewall_enabled, instance_name_prefix, networks, metadata, userdata, tags)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

	metaDataJson, err := convertToJson(template.MetaData)
	if err != nil {
		log.Fatal().Err(err)
	}

	tagsJson, err := convertToJson(template.Tags)
	if err != nil {
		log.Fatal().Err(err)
	}

	networksList := strings.Join(template.Networks, ",")

	_, err = db.ExecEx(ctx, sqlStatement, nil,
		template.TemplateName, template.Package, template.ImageID,
		accountID, template.FirewallEnabled, template.InstanceNamePrefix, networksList, metaDataJson,
		template.UserData, tagsJson)
	if err != nil {
		panic(err)
	}
}

func RemoveTemplate(ctx context.Context, identifier int64, accountID int) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return
	}

	sqlStatement := `UPDATE triton.tsg_templates
SET archived = true
WHERE id = $1 and account_id = $2`

	_, err := db.ExecEx(ctx, sqlStatement, nil, identifier, accountID)
	if err != nil {
		panic(err)
	}
}

func convertToJson(data map[string]string) (string, error) {
	if data == nil {
		return "", nil
	}

	log.Debug().Msg("Found data")
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
