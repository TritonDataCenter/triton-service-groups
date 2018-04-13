//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package groups_v1

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"github.com/joyent/triton-service-groups/server/handlers"
	"github.com/rs/zerolog/log"
)

func FindGroups(ctx context.Context, accountID string) ([]*ServiceGroup, error) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return nil, handlers.ErrNoConnPool
	}

	var groups []*ServiceGroup

	sqlStatement := `
SELECT id, name, account_id, template_id, capacity, created_at, updated_at
FROM tsg_groups
WHERE account_id = $1
AND archived = false;`

	rows, err := db.QueryEx(ctx, sqlStatement, nil, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			group     ServiceGroup
			createdAt pgtype.Timestamp
			updatedAt pgtype.Timestamp
		)

		err := rows.Scan(
			&group.ID,
			&group.GroupName,
			&group.AccountID,
			&group.TemplateID,
			&group.Capacity,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		group.CreatedAt = createdAt.Time
		group.UpdatedAt = updatedAt.Time

		groups = append(groups, &group)
	}

	return groups, nil
}

func FindGroupByID(ctx context.Context, key string, accountID string) (*ServiceGroup, bool) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return nil, false
	}

	var (
		group     ServiceGroup
		createdAt pgtype.Timestamp
		updatedAt pgtype.Timestamp
	)

	sqlStatement := `
SELECT id, name, account_id, template_id, capacity, created_at, updated_at
FROM tsg_groups
WHERE account_id = $2 and id = $1
AND archived = false
`

	err := db.QueryRowEx(ctx, sqlStatement, nil, key, accountID).Scan(
		&group.ID,
		&group.GroupName,
		&group.AccountID,
		&group.TemplateID,
		&group.Capacity,
		&createdAt,
		&updatedAt,
	)
	switch err {
	case nil:
		group.CreatedAt = createdAt.Time
		group.UpdatedAt = updatedAt.Time

		return &group, true
	case pgx.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil, false
	default:
		panic(err)
	}
}

func FindGroupByName(ctx context.Context, name string, accountID string) (*ServiceGroup, bool) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return nil, false
	}

	var (
		group     ServiceGroup
		createdAt pgtype.Timestamp
		updatedAt pgtype.Timestamp
	)

	sqlStatement := `
SELECT id, name, account_id, template_id, capacity, created_at, updated_at
FROM tsg_groups
WHERE account_id = $2 and name = $1
AND archived = false;
`
	err := db.QueryRowEx(ctx, sqlStatement, nil, name, accountID).Scan(
		&group.ID,
		&group.GroupName,
		&group.AccountID,
		&group.TemplateID,
		&group.Capacity,
		&createdAt,
		&updatedAt,
	)
	switch err {
	case nil:
		group.CreatedAt = createdAt.Time
		group.UpdatedAt = updatedAt.Time

		return &group, true
	case pgx.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil, false
	default:
		panic(err)
	}
}

func SaveGroup(ctx context.Context, accountID string, group *ServiceGroup) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return
	}

	sqlStatement := `
INSERT INTO tsg_groups (name, template_id, capacity, account_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
`
	_, err := db.ExecEx(ctx, sqlStatement, nil,
		group.GroupName,
		group.TemplateID,
		group.Capacity,
		accountID,
	)
	if err != nil {
		panic(err)
	}
}

func UpdateGroup(ctx context.Context, uuid string, accountID string, group *ServiceGroup) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return
	}

	sqlStatement := `
UPDATE tsg_groups
SET template_id = $3, capacity = $4, updated_at = NOW()
WHERE id = $1 and account_id = $2
`
	_, err := db.ExecEx(ctx, sqlStatement, nil,
		uuid,
		accountID,
		group.TemplateID,
		group.Capacity,
	)
	if err != nil {
		panic(err)
	}
}

func RemoveGroup(ctx context.Context, identifier string, accountID string) {
	db, ok := handlers.GetDBPool(ctx)
	if !ok {
		log.Fatal().Err(handlers.ErrNoConnPool)
		return
	}

	sqlStatement := `
UPDATE tsg_groups
SET archived = true, updated_at = NOW()
WHERE id = $1 and account_id = $2
`
	_, err := db.ExecEx(ctx, sqlStatement, nil, identifier, accountID)
	if err != nil {
		panic(err)
	}
}
