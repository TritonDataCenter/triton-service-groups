package templates_v1

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/jackc/pgx"
)

func FindTemplateBy(db *pgx.ConnPool, key string, accountId string) (*MachineTemplate, bool) {
	var template MachineTemplate

	sqlStatement := `SELECT name, package, image_id, firewall_enabled, metadata, userdata, tags  
FROM triton.tsg_templates 
WHERE name = $1 and account_id = $2;`

	var userDataJson string
	var metaDataJson string
	var tagsJson string

	err := db.QueryRowEx(context.TODO(), sqlStatement, nil, key, accountId).
		Scan(&template.Name,
			&template.Package,
			&template.ImageID,
			&template.FirewallEnabled,
			&metaDataJson,
			&userDataJson,
			&tagsJson)
	switch err {
	case nil:
		err = json.Unmarshal([]byte(userDataJson), template.UserData)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(metaDataJson), template.MetaData)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(tagsJson), template.Tags)
		if err != nil {
			panic(err)
		}

		return &template, true
	case pgx.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil, false
	default:
		panic(err)
	}
}

func FindTemplates(db *pgx.ConnPool, accountId string) ([]*MachineTemplate, error) {
	var templates []*MachineTemplate

	sqlStatement := `SELECT name, package, image_id, firewall_enabled, metadata, userdata, tags 
FROM triton.tsg_templates
WHERE account_id = $1;`

	var userDataJson string
	var metaDataJson string
	var tagsJson string

	rows, err := db.QueryEx(context.TODO(), sqlStatement, nil, accountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var template MachineTemplate
		err := rows.Scan(&template.Name,
			&template.Package,
			&template.ImageID,
			&template.FirewallEnabled,
			&metaDataJson,
			&userDataJson,
			&tagsJson)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(userDataJson), template.UserData)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(metaDataJson), template.MetaData)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(tagsJson), template.Tags)
		if err != nil {
			panic(err)
		}

		templates = append(templates, &template)
	}

	return templates, nil
}

func SaveTemplate(db *pgx.ConnPool, accountId string, template *MachineTemplate) {
	sqlStatement := `
INSERT INTO triton.tsg_templates (name, package, image_id, account_id, firewall_enabled, metadata, userdata, tags) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

	metaDataJson, _ := json.Marshal(template.MetaData)
	userDataJson, _ := json.Marshal(template.UserData)
	tagsJson, _ := json.Marshal(template.Tags)

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil,
		template.Name, template.Package, template.ImageID,
		accountId, template.FirewallEnabled, metaDataJson,
		userDataJson, tagsJson)
	if err != nil {
		panic(err)
	}
}

func UpdateTemplate(db *pgx.ConnPool, name string, accountId string, template *MachineTemplate) {
	sqlStatement := `
Update triton.tsg_templates 
SET package = $3, image_id = $4, firewall_enabled = $5, metadata = $6, userdata = $7, tags = $8
WHERE name = $1 and account_id = $2
`

	metaDataJson, _ := json.Marshal(template.MetaData)
	userDataJson, _ := json.Marshal(template.UserData)
	tagsJson, _ := json.Marshal(template.Tags)

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil,
		name, accountId, template.Package, template.ImageID, template.FirewallEnabled,
		metaDataJson, userDataJson, tagsJson)
	if err != nil {
		panic(err)
	}
}

func RemoveTemplate(db *pgx.ConnPool, name string, accountId string) {
	sqlStatement := `DELETE FROM triton.tsg_templates WHERE name = $1 and account_id = $2`

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil, name, accountId)
	if err != nil {
		panic(err)
	}
}
