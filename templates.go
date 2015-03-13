package main

import "text/template"

/* Template helper functions */
var fns = template.FuncMap{
	"plus1": func(x int) int {
		return x + 1
	},
}

/* Tables */

const TABLE_TEMPLATE = `package {{.Options.PackageName}}

/* *********************************************************** */
/* This file was automatically generated by pgtogogen.         */
/* Do not modify this file unless you know what you are doing. */
/* *********************************************************** */

import (
	{{range $key, $value := .GoTypesToImport}}"{{$value}}"
	{{end}}	
)

const {{.GoFriendlyName}}_DB_TABLE_NAME string = "{{.TableName}}"

type {{.GoFriendlyName}} struct {
	{{range .Columns}}{{.GoName}} {{.GoType}} // IsPK: {{.IsPK}} , IsCompositePK: {{.IsCompositePK}}, IsFK: {{.IsFK}}
	{{end}}	
}`

const TABLE_TEMPLATE_CUSTOM = `package {{.Options.PackageName}}

/* *********************************************************** **/
/* This file is generated by pgtogogen FIRST-TIME ONLY.         */
/* It will not subsequently overwrite it if it already exists.  */
/* Use this file to create your custom extension functionality. */
/* ************************************************************ */

/*
import (
	{{range $key, $value := .GoTypesToImport}}"{{$value}}"
	{{end}}	
)
*/

`

/* Columns */

const PK_GETTER_TEMPLATE = `{{$colCount := len .ParentTable.Columns}}
// Queries the database for a single row based on the specified {{.GoName}} value.
// Returns a pointer to a {{.ParentTable.GoFriendlyName}} structure if a record was found,
// otherwise it returns nil.
func {{.ParentTable.GoFriendlyName}}GetBy{{.GoName}}(inputParam{{.GoName}} {{.GoType}}) (returnStruct *{{.ParentTable.GoFriendlyName}}, err error) {
	
	returnStruct = nil
	err = nil
	
	var errorPrefix = "{{.ParentTable.GoFriendlyName}}GetBy{{.GoName}}() ERROR: "
	
	currentDbHandle := GetDb()
	if currentDbHandle == nil {
		return nil, NewModelsErrorLocal(errorPrefix, "the database handle is nil")
	}

	// define receiving params for the row iteration
	{{range .ParentTable.Columns}}var param{{.GoName}} {{.GoType}}
	{{end}}

	// define the select query
	var query = "{{.ParentTable.GenericSelectQuery}} WHERE {{.Name}} = $1";

	// we are aiming for a single row so we will use Query Row	
	err = currentDbHandle.QueryRow(query, inputParam{{.GoName}}).Scan({{range $i, $e := .ParentTable.Columns}}&param{{$e.GoName}}{{if ne (plus1 $i) $colCount}},{{end}}{{end}})
    switch {
    case err == sql.ErrNoRows:
            // no such row found, return nil and nil
			return nil, nil
    case err != nil:
            return nil, NewModelsError(errorPrefix + "fatal error running the query:", err)
    default:
           	// create the return structure as a pointer of the type
			returnStruct = &{{.ParentTable.GoFriendlyName}}{
				{{range .ParentTable.Columns}}{{.GoName}}: param{{.GoName}},
				{{end}}
			}
			// return the structure
			return returnStruct, nil
    }			
}
`

const PK_SELECT_TEMPLATE = `{{$colCount := len .ParentTable.Columns}}
func {{.ParentTable.GoFriendlyName}}GetBy{{.GoName}}(inputParam{{.GoName}} {{.GoType}}) (returnStruct *{{.ParentTable.GoFriendlyName}}, err error) {
	
	returnStruct = nil
	err = nil
	
	var errorPrefix = "{{.ParentTable.GoFriendlyName}}GetBy{{.GoName}}() ERROR: "
	
	currentDbHandle := GetDb()
	if currentDbHandle == nil {
		return nil, NewModelsErrorLocal(errorPrefix, "the database handle is nil")
	}

	// define receiving params for the row iteration
	{{range .ParentTable.Columns}}var param{{.GoName}} {{.GoType}}
	{{end}}

	// define the select query
	var query = "{{.ParentTable.GenericSelectQuery}} FROM {{.ParentTable.TableName}} WHERE {{.Name}} = $1";

	rows, err := currentDbHandle.Query(query, inputParam{{.GoName}})

	if err != nil {
		return nil, NewModelsError(errorPrefix + "fatal error running the query:", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan({{range $i, $e := .ParentTable.Columns}}&param{{$e.GoName}}{{if ne (plus1 $i) $colCount}},{{end}}{{end}})
		if err != nil {
			return nil, NewModelsError(errorPrefix + "fatal error scanning the fields in the current row:",err.Error)
		}		

		// create the return structure as a pointer of the type
		returnStruct = &{{.ParentTable.GoFriendlyName}}{
			{{range .ParentTable.Columns}}{{.GoName}}: param{{.GoName}},
			{{end}}
		}		

	}
	err = rows.Err()
	if err != nil {
		return nil, NewModelsError(errorPrefix + "fatal generic rows error:", err.Error)
	}
	
	// return the structure
	return returnStruct, err	
}
`
