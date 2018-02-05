//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package groups_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx"
)

func FindGroups(db *pgx.ConnPool, accountId string) ([]*ServiceGroup, error) {
	var groups []*ServiceGroup

	sqlStatement := `SELECT name, template_id, capacity, datacenter, health_check_interval, COALESCE(instance_tags, '') 
FROM triton.tsg_groups
WHERE account_id = $1 
AND archived = false;`

	var instanceTagsJson string
	var datacenterList string

	rows, err := db.QueryEx(context.TODO(), sqlStatement, nil, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var group ServiceGroup
		err := rows.Scan(&group.GroupName,
			&group.TemplateId,
			&group.Capacity,
			&datacenterList,
			&group.HealthCheckInterval,
			&instanceTagsJson)
		if err != nil {
			return nil, err
		}

		tags, err := convertFromJson(instanceTagsJson)
		if err != nil {
			panic(err)
		}
		group.InstanceTags = tags

		group.DataCenter = strings.Split(datacenterList, ",")

		groups = append(groups, &group)
	}

	return groups, nil
}

func FindGroupBy(db *pgx.ConnPool, key string, accountId string) (*ServiceGroup, bool) {
	var group ServiceGroup

	sqlStatement := `SELECT name, template_id, capacity, datacenter, health_check_interval, COALESCE(instance_tags, '') 
FROM triton.tsg_groups
WHERE account_id = $1 and name = $2
AND archived = false;`

	var instanceTagsJson string
	var datacenterList string

	err := db.QueryRowEx(context.TODO(), sqlStatement, nil, key, accountId).
		Scan(&group.GroupName,
			&group.TemplateId,
			&group.Capacity,
			&datacenterList,
			&group.HealthCheckInterval,
			&instanceTagsJson)
	switch err {
	case nil:
		instanceTags, err := convertFromJson(instanceTagsJson)
		if err != nil {
			panic(err)
		}
		group.InstanceTags = instanceTags

		group.DataCenter = strings.Split(datacenterList, ",")

		return &group, true
	case pgx.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil, false
	default:
		panic(err)
	}
}

func SaveGroup(db *pgx.ConnPool, accountId string, group *ServiceGroup) {
	sqlStatement := `
INSERT INTO triton.tsg_groups (name, template_id, capacity, account_id, datacenter, health_check_interval, instance_tags) 
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

	tagsJson, err := convertToJson(group.InstanceTags)
	if err != nil {
		log.Fatal(err)
	}

	datacenterList := strings.Join(group.DataCenter, ",")

	_, err = db.ExecEx(context.TODO(), sqlStatement, nil,
		group.GroupName, group.TemplateId, group.Capacity,
		accountId, datacenterList, group.HealthCheckInterval,
		tagsJson)
	if err != nil {
		panic(err)
	}
}

func UpdateGroup(db *pgx.ConnPool, name string, accountId string, group *ServiceGroup) {
	sqlStatement := `
Update triton.tsg_groups 
SET template_id = $3, capacity = $4, datacenter = $5, health_check_interval = $6, instance_tags = $7
WHERE name = $1 and account_id = $2
`

	tagsJson, err := convertToJson(group.InstanceTags)
	if err != nil {
		log.Fatal(err)
	}

	datacenterList := strings.Join(group.DataCenter, ",")

	_, err = db.ExecEx(context.TODO(), sqlStatement, nil,
		name, accountId, group.TemplateId, group.Capacity, datacenterList,
		group.HealthCheckInterval, tagsJson)
	if err != nil {
		panic(err)
	}
}

func RemoveGroup(db *pgx.ConnPool, name string, accountId string) {
	sqlStatement := `UPDATE triton.tsg_groups 
SET archived = true 
WHERE name = $1 and account_id = $2`

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil, name, accountId)
	if err != nil {
		panic(err)
	}
}

func convertToJson(data map[string]string) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	log.Printf("Found data")
	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return json, nil
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
