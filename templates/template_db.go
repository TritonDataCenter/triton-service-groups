package templates_v1

import (
	"fmt"

	"context"

	"github.com/jackc/pgx"
)

func dbConnection() (*pgx.Conn, error) {
	config := pgx.ConnConfig{
		Host:     "localhost",
		Database: "tsg",
		Port:     26257,
		User:     "root",
	}

	return pgx.Connect(config)
}

func FindTemplateBy(key string) (*MachineTemplate, bool) {
	var template MachineTemplate

	sqlStatement := `SELECT name, package, image_id FROM triton.tsg_templates WHERE name = $1;`

	db, err := dbConnection()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.QueryRowEx(context.TODO(), sqlStatement, nil, key).Scan(&template.Name, &template.Package, &template.ImageID)
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

func FindTemplates() ([]*MachineTemplate, error) {
	var templates []*MachineTemplate

	sqlStatement := `SELECT name, package, image_id FROM triton.tsg_templates;`

	db, err := dbConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()

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

func SaveTemplate(template *MachineTemplate) {
	sqlStatement := `
INSERT INTO triton.tsg_templates (name, package, image_id) 
VALUES ($1, $2, $3)
`

	db, err := dbConnection()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.ExecEx(context.TODO(), sqlStatement, nil, template.Name, template.Package, template.ImageID)
	if err != nil {
		panic(err)
	}
}

func UpdateTemplate(name string, template *MachineTemplate) {
	sqlStatement := `
Update triton.tsg_templates 
SET package = $2, image_id = $3
WHERE name = $1
`

	db, err := dbConnection()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.ExecEx(context.TODO(), sqlStatement, nil, name, template.Package, template.ImageID)
	if err != nil {
		panic(err)
	}
}

func RemoveTemplate(name string) {
	sqlStatement := `DELETE FROM triton.tsg_templates WHERE name = $1`

	db, err := dbConnection()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	_, err = db.ExecEx(context.TODO(), sqlStatement, nil, name)
	if err != nil {
		panic(err)
	}
}
