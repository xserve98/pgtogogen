package main

/* Views */

const VIEW_TEMPLATE = `package {{.Options.PackageName}}

/* *********************************************************** */
/* This file was automatically generated by pgtogogen.         */
/* Do not modify this file unless you know what you are doing. */
/* *********************************************************** */

import (	
	"bytes"
	"sync"
	{{range $key, $value := .GoTypesToImport}}"{{$value}}"
	{{end}}
)

const {{.GoFriendlyName}}_DB_VIEW_NAME string = "{{.DbName}}"

type {{.GoFriendlyName}} struct {
	{{range .Columns}}{{.GoName}} {{.GoType}} 
	{{end}}		
}

// fake, interal type to allow a singleton structure that would hold static-like methods
type t{{.GoFriendlyName}}Utils struct {
		
}

{{if .IsMaterialized}}
{{$functionName := "RefreshMaterializedView"}}{{$sourceStructName := print "source" .GoFriendlyName}}
// Refreshes the materialized view and updates it with the latest data from the underlying 
// data entities
func (utilRef *t{{.GoFriendlyName}}Utils) {{$functionName}}()  error {
						
	var errorPrefix = "{{.GoFriendlyName}}Utils.{{$functionName}}() ERROR: "
	
	currentDbHandle := GetDb()
	if currentDbHandle == nil {
		return NewModelsErrorLocal(errorPrefix, "the database handle is nil")
	}
	
	_, err := currentDbHandle.Exec("REFRESH MATERIALIZED VIEW {{.DbName}};")
	if err != nil {
		return NewModelsError(errorPrefix + "currentDbHandle.Exec error:",err)
	}
	
	return  nil	
	
}
{{end}}
`

const VIEW_TEMPLATE_CUSTOM = `package {{.Options.PackageName}}

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
