package main

/* Views */

const VIEW_TEMPLATE = `package {{.Options.PackageName}}

/* *********************************************************** */
/* This file was automatically generated by pgtogogen.         */
/* Do not modify this file unless you know what you are doing. */
/* *********************************************************** */

import (	
	"sync"
	pgx "{{.Options.PgxImport}}"
	pgtype "{{.Options.PgTypeImport}}"
	{{range $key, $value := .GoTypesToImport}}"{{$value}}"
	{{end}}
)


// this is a dummy variable, just to use the pgx package
var pgxPackageView{{.GoFriendlyName}}ReferenceErrDeadConn = pgx.ErrDeadConn

const {{.GoFriendlyName}}_DB_VIEW_NAME string = "{{.DbName}}"

type {{.GoFriendlyName}} struct {
	{{range .Columns}}// database field name: {{.DbName}}
	{{.GoName}} {{.GoType}}
	{{if .Nullable}}{{.GoName}}_IsNotNull bool // if true, it means the value is not null
	{{end}}
	{{end}}		
}

{{ $tableGoName := .GoFriendlyName}}
/* Sorting helper containers */
{{range $i, $e := .Columns}}
// By{{$e.GoName}} implements sort.Interface for []{{$tableGoName}} based on
// the {{$e.GoName}} field. Usage: sort.Sort(Sort{{$tableGoName}}By{{$e.GoName}}(anyGiven{{$tableGoName}}Slice))
type Sort{{$tableGoName}}By{{$e.GoName}} []{{$tableGoName}}

func (a Sort{{$tableGoName}}By{{$e.GoName}}) Len() int           { return len(a) }
func (a Sort{{$tableGoName}}By{{$e.GoName}}) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Sort{{$tableGoName}}By{{$e.GoName}}) Less(i, j int) bool { return LessComparatorFor_{{$e.GoType}}(a[i].{{$e.GoName}},a[j].{{$e.GoName}}) }
{{end}}
{{range .Columns}}func (t *{{$tableGoName}}) Set{{.GoName}}(val {{.GoType}} {{if .Nullable}}, notNull bool{{end}}) {
	t.{{.GoName}} = val
	{{if .Nullable}}t.{{.GoName}}_IsNotNull = notNull{{end}}
}
{{end}}

// fake, internal type to allow a singleton structure that would hold static-like methods
type t{{.GoFriendlyName}}Utils struct {
		
	// instance of a CacheFor{{.GoFriendlyName}} structure
	Cache CacheFor{{.GoFriendlyName}}		
		
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

{{$functionNameConc := "RefreshMaterializedViewConcurrently"}}{{$sourceStructName := print "source" .GoFriendlyName}}
// Refreshes the materialized view concurrently, and updates it with the latest data from 
// the underlying data entities. A concurrent refresh means that the view is accessible to reading
// by other threads, but it may take longer than the non-concurrent operation. 
// This refresh mode is only available in Postgres versions 9.4 and higher and it will fail unless at least one
// unique index, without a WHERE clause is defined on the view
func (utilRef *t{{.GoFriendlyName}}Utils) {{$functionNameConc}}()  error {
						
	var errorPrefix = "{{.GoFriendlyName}}Utils.{{$functionNameConc}}() ERROR: "
	
	currentDbHandle := GetDb()
	if currentDbHandle == nil {
		return NewModelsErrorLocal(errorPrefix, "the database handle is nil")
	}
	
	_, err := currentDbHandle.Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY {{.DbName}};")
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
