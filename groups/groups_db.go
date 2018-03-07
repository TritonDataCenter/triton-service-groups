//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package groups_v1

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
)

func FindGroups(db *pgx.ConnPool, accountId string) ([]*ServiceGroup, error) {
	var groups []*ServiceGroup

	sqlStatement := `SELECT id, name, account_id, template_id, capacity, health_check_interval 
FROM triton.tsg_groups
WHERE account_id = $1 
AND archived = false;`

	rows, err := db.QueryEx(context.TODO(), sqlStatement, nil, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var group ServiceGroup
		err := rows.Scan(&group.ID,
			&group.GroupName,
			&group.AccountID,
			&group.TemplateID,
			&group.Capacity,
			&group.HealthCheckInterval)
		if err != nil {
			return nil, err
		}

		groups = append(groups, &group)
	}

	return groups, nil
}

func FindGroupBy(db *pgx.ConnPool, key string, accountId string) (*ServiceGroup, bool) {
	var group ServiceGroup

	sqlStatement := `SELECT id, name, account_id, template_id, capacity, health_check_interval 
FROM triton.tsg_groups
WHERE account_id = $2 and name = $1
AND archived = false;`

	err := db.QueryRowEx(context.TODO(), sqlStatement, nil, key, accountId).
		Scan(&group.ID,
			&group.GroupName,
			&group.AccountID,
			&group.TemplateID,
			&group.Capacity,
			&group.HealthCheckInterval)
	switch err {
	case nil:
		return &group, true
	case pgx.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil, false
	default:
		panic(err)
	}
}

func FindGroupByID(db *pgx.ConnPool, key int64, accountId string) (*ServiceGroup, bool) {
	var group ServiceGroup

	sqlStatement := `SELECT id, name, account_id, template_id, capacity, health_check_interval 
FROM triton.tsg_groups
WHERE account_id = $2 and id = $1
AND archived = false;`

	err := db.QueryRowEx(context.TODO(), sqlStatement, nil, key, accountId).
		Scan(&group.ID,
			&group.GroupName,
			&group.AccountID,
			&group.TemplateID,
			&group.Capacity,
			&group.HealthCheckInterval)
	switch err {
	case nil:
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
INSERT INTO triton.tsg_groups (name, template_id, capacity, account_id, health_check_interval) 
VALUES ($1, $2, $3, $4, $5)
`
	_, err := db.ExecEx(context.TODO(), sqlStatement, nil,
		group.GroupName, group.TemplateID, group.Capacity,
		accountId, group.HealthCheckInterval)
	if err != nil {
		panic(err)
	}
}

func UpdateGroup(db *pgx.ConnPool, name string, accountId string, group *ServiceGroup) {
	sqlStatement := `
Update triton.tsg_groups 
SET template_id = $3, capacity = $4, health_check_interval = $5
WHERE name = $1 and account_id = $2
`

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil,
		name, accountId, group.TemplateID, group.Capacity, group.HealthCheckInterval)
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
