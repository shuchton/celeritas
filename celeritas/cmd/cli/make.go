package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3 string) error {
	switch arg2 {
	case "key":
		rnd := cel.RandomString(32)
		color.Yellow("32 character encryptions key: %s", rnd)
	case "migration":
		dbType := cel.DB.DataBaseType
		if arg3 == "" {
			exitGracefully(errors.New("migration must have a name"))
		}

		filename := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := cel.RootPath + "/migrations/" + filename + "." + dbType + ".up.sql"
		downFile := cel.RootPath + "/migrations/" + filename + "." + dbType + ".down.sql"

		err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}
	case "auth":
		err := doAuth()
		if err != nil {
			return err
		}
	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("handler must have a name"))
		}

		fileName := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}
	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("model must have a name"))
		}
		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		model := string(data)

		plur := pluralize.NewClient()

		var modelName = arg3
		var tableName = arg3
		if plur.IsPlural(arg3) {
			modelName = plur.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(plur.Plural(arg3))
		}

		fileName := cel.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + " already exists"))
		}

		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

		err = copyDataToFile([]byte(model), fileName)

		if err != nil {
			exitGracefully(err)
		}

	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}
	}
	return nil
}
