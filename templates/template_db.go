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
	"github.com/jackc/pgx/pgtype"
	"github.com/joyent/triton-service-groups/convert"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/rs/zerolog/log"
)

func CheckTemplateExistsByName(ctx context.Context, templateName, accountID string) (bool, error) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		return false, handlers.ErrNoConnPool
	}

	var exists bool

	sql := `
SELECT EXISTS
  (SELECT 1
   FROM tsg_templates
   WHERE (template_name = $1
          AND account_id = $2)
     AND archived IS FALSE);`

	err := db.QueryRowEx(ctx, sql, nil, templateName, accountID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func CheckTemplateAllocationByID(ctx context.Context, templateID, accountID string) (bool, error) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		return false, handlers.ErrNoConnPool
	}

	var allocated bool

	sql := `
SELECT EXISTS
  (SELECT 1
   FROM tsg_templates AS t,
        tsg_groups AS g
   WHERE (t.id = g.template_id
          AND t.account_id = g.account_id)
     AND (t.id = $1
          AND t.account_id = $2)
     AND g.archived IS FALSE);`

	err := db.QueryRowEx(ctx, sql, nil, templateID, accountID).Scan(&allocated)
	if err != nil {
		return false, err
	}

	return allocated, nil
}

func FindTemplateByName(ctx context.Context, key string, accountID string) (*InstanceTemplate, bool) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		return nil, false
	}

	sqlStatement := `
SELECT id, template_name, package, image_id, firewall_enabled, networks, COALESCE(metadata,''), userdata, COALESCE(tags,''), created_at
FROM tsg_templates
WHERE template_name = $1 and account_id = $2
AND archived = false
`

	var (
		template     InstanceTemplate
		metaDataJson string
		tagsJson     string
		networksList string
		templateID   pgtype.UUID
		createdAt    pgtype.Timestamp
	)

	err := db.QueryRowEx(ctx, sqlStatement, nil, key, accountID).Scan(
		&templateID,
		&template.TemplateName,
		&template.Package,
		&template.ImageID,
		&template.FirewallEnabled,
		&networksList,
		&metaDataJson,
		&template.UserData,
		&tagsJson,
		&createdAt,
	)
	switch err {
	case nil:
		template.ID = convert.BytesToUUID(templateID.Bytes)

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

		template.CreatedAt = createdAt.Time

		return &template, true
	case pgx.ErrNoRows:
		return nil, false
	default:
		return nil, false
	}
}

func FindTemplateByID(ctx context.Context, key string, accountID string) (*InstanceTemplate, bool) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {

		return nil, false
	}

	sqlStatement := `
SELECT id, template_name, package, image_id, firewall_enabled, networks, COALESCE(metadata,''), userdata, COALESCE(tags,''), created_at
FROM tsg_templates
WHERE id = $1 and account_id = $2
AND archived = false
`

	var (
		template     InstanceTemplate
		metaDataJson string
		tagsJson     string
		networksList string
		templateID   pgtype.UUID
		createdAt    pgtype.Timestamp
	)

	err := db.QueryRowEx(ctx, sqlStatement, nil, key, accountID).Scan(
		&templateID,
		&template.TemplateName,
		&template.Package,
		&template.ImageID,
		&template.FirewallEnabled,
		&networksList,
		&metaDataJson,
		&template.UserData,
		&tagsJson,
		&createdAt,
	)
	switch err {
	case nil:
		template.ID = convert.BytesToUUID(templateID.Bytes)

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

		template.CreatedAt = createdAt.Time

		return &template, true
	case pgx.ErrNoRows:
		return nil, false
	default:
		return nil, false
	}
}

func FindTemplates(ctx context.Context, accountID string) ([]*InstanceTemplate, error) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		return nil, handlers.ErrNoConnPool
	}

	sqlStatement := `SELECT id, template_name, package, image_id, firewall_enabled, networks, COALESCE(metadata,''), userdata, COALESCE(tags, ''), created_at
FROM tsg_templates
WHERE account_id = $1
AND archived = false;`

	var (
		templates    []*InstanceTemplate
		metaDataJson string
		tagsJson     string
		networksList string
		templateID   pgtype.UUID
		createdAt    pgtype.Timestamp
	)

	rows, err := db.QueryEx(ctx, sqlStatement, nil, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var template InstanceTemplate
		err := rows.Scan(
			&templateID,
			&template.TemplateName,
			&template.Package,
			&template.ImageID,
			&template.FirewallEnabled,
			&networksList,
			&metaDataJson,
			&template.UserData,
			&tagsJson,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		template.ID = convert.BytesToUUID(templateID.Bytes)

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

		template.CreatedAt = createdAt.Time

		templates = append(templates, &template)
	}

	return templates, nil
}

func SaveTemplate(ctx context.Context, accountID string, template *InstanceTemplate) error {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		return handlers.ErrNoConnPool
	}

	sqlStatement := `
INSERT INTO tsg_templates (template_name, package, image_id, account_id, firewall_enabled, networks, metadata, userdata, tags, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
`

	metaDataJson, err := convertToJson(template.MetaData)
	if err != nil {
		return err
	}

	tagsJson, err := convertToJson(template.Tags)
	if err != nil {
		return err
	}

	networksList := strings.Join(template.Networks, ",")

	_, err = db.ExecEx(ctx, sqlStatement, nil,
		template.TemplateName,
		template.Package,
		template.ImageID,
		accountID,
		template.FirewallEnabled,
		networksList,
		metaDataJson,
		template.UserData,
		tagsJson,
	)
	if err != nil {
		return err
	}

	return nil
}

func RemoveTemplate(ctx context.Context, identifier string, accountID string) error {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		return handlers.ErrNoConnPool
	}

	sqlStatement := `UPDATE triton.tsg_templates
SET archived = true
WHERE id = $1 and account_id = $2`

	_, err := db.ExecEx(ctx, sqlStatement, nil, identifier, accountID)
	if err != nil {
		return err
	}

	return nil
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
