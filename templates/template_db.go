package templates_v1

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
)

func FindTemplateBy(db *pgx.ConnPool, key string) (*MachineTemplate, bool) {
	var template MachineTemplate

	sqlStatement := `SELECT name, package, image_id FROM triton.tsg_templates WHERE name = $1;`

	err := db.QueryRowEx(context.TODO(), sqlStatement, nil, key).
		Scan(&template.Name, &template.Package, &template.ImageID)
	switch err {
	case nil:
		return &template, true
	case pgx.ErrNoRows:
		fmt.Println("No rows were returned!")
		return nil, false
	default:
		panic(err)
	}
}

func FindTemplates(db *pgx.ConnPool) ([]*MachineTemplate, error) {
	var templates []*MachineTemplate

	sqlStatement := `SELECT name, package, image_id FROM triton.tsg_templates;`

	rows, err := db.QueryEx(context.TODO(), sqlStatement, nil)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var template MachineTemplate
		err := rows.Scan(&template.Name, &template.Package, &template.ImageID)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &template)
	}

	return templates, nil
}

func SaveTemplate(db *pgx.ConnPool, template *MachineTemplate) {
	sqlStatement := `
INSERT INTO triton.tsg_templates (name, package, image_id) 
VALUES ($1, $2, $3)
`

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil, template.Name, template.Package, template.ImageID)
	if err != nil {
		panic(err)
	}
}

func UpdateTemplate(db *pgx.ConnPool, name string, template *MachineTemplate) {
	sqlStatement := `
Update triton.tsg_templates 
SET package = $2, image_id = $3
WHERE name = $1
`

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil, name, template.Package, template.ImageID)
	if err != nil {
		panic(err)
	}
}

func RemoveTemplate(db *pgx.ConnPool, name string) {
	sqlStatement := `DELETE FROM triton.tsg_templates WHERE name = $1`

	_, err := db.ExecEx(context.TODO(), sqlStatement, nil, name)
	if err != nil {
		panic(err)
	}
}
