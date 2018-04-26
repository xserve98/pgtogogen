package main

const FUNCTION_TEMPLATE_PREFIX = `package {{.Options.PackageName}}

/* *********************************************************** */
/* This file was automatically generated by pgtogogen.         */
/* Do not modify this file unless you know what you are doing. */
/* *********************************************************** */

import (
	pgx "{{.Options.PgxImport}}"
	pgtype "{{.Options.PgTypeImport}}"
	{{range $key, $value := .GoTypesToImport}}"{{$value}}"
	{{end}}	
)

// this is a dummy variable, just to use the pgx package
var pgxErrDeadConnFunc = pgx.ErrDeadConn

// this is a dummy variable, just to use the pgtypes package
const pgtypesDummyFnPlaceholder = pgtype.Present

// Utility-oriented, internal type to allow a singleton structure that would hold static-like methods
// and global, single-instance settings
type tFunctionUtils struct {}

var Functions tFunctionUtils

`

const COMMON_CODE_FUNCTION_QUERY = `rows, err := currentDbHandle.Query(JoinStringParts(queryParts,""), {{range $i, $e := .Parameters}}param{{.GoFriendlyName}}{{if ne (plus1 $i) $paramCount}},{{end}} {{end}})	
	
	if err != nil {
		return {{if not .IsReturnVoid}}returnVal,{{end}} NewModelsError(errorPrefix + " fatal error running the function statement:", err)
	}
	defer rows.Close()

	{{if .IsReturnVoid}}
	// this function returns void, do not attempt doing anything else
	return nil
	{{else}}//BEGIN: non void operations

	{{if .IsReturnUserDefined}}// BEGIN: if any nullable fields, create temporary nullable variables to receive null values
	{{range $i, $e := .Columns}}{{if .Nullable}}var nullable{{$e.GoName}} {{$e.GoNullableType}} 
	{{end}}{{end}}
	// END: if any nullable fields, create temporary nullable variables to receive null values{{end}}

	for rows.Next() {

		// create a new variable of type {{.ReturnGoType}}  {{$instanceVarName := print "current" .ReturnGoType}}			
		var {{$instanceVarName}} {{.ReturnGoType}}		
		
	{{if .IsReturnUserDefined}}		// BEGIN: User-defined return type collection (slice)

		{{$colCount := len .Columns}}
		err := rows.Scan({{range $i, $e := .Columns}}{{if .Nullable}}&nullable{{$e.GoName}}{{else}}&{{$instanceVarName}}.{{$e.GoName}}{{end}}{{if ne (plus1 $i) $colCount}},{{end}}{{end}})
		if err != nil {
			return returnVal, NewModelsError(errorPrefix + " error during rows.Scan():", err)
		}
		
		// BEGIN: assign any nullable values to the nullable fields inside the struct appropriately
		{{range $i, $e := .Columns}}{{if .Nullable}} {{$instanceVarName}}.Set{{.GoName}}(nullable{{$e.GoName}}.GetValue(), nullable{{$e.GoName}}.Valid)
		{{end}}{{end}}
		// END: assign any nullable values to the nullable fields inside the struct appropriately	
		
		// END: User-defined (internal types) return type collection (slice)		
	{{else}}	// BEGIN: Not-user-defined (internal types) return type collection (slice)
	
		err = rows.Scan(&{{$instanceVarName}})
		if err != nil {
			return returnVal, NewModelsError(errorPrefix + " error during rows.Scan():", err)
		}							
		// END: Not-User-defined (internal types) return type collection (slice)			
	{{end}}
				
		// a set is returned (expect one or more records)
		returnVal = append(returnVal, {{$instanceVarName}})

	}
	err = rows.Err()
	if err != nil {
		return returnVal, NewModelsError(errorPrefix + " error during rows.Next() iterations:", err)
	}
	
	{{end}} //END: non void operations
	
	return {{if not .IsReturnVoid}}returnVal,{{end}} nil
`

const COMMON_CODE_FUNCTION_QUERYROW = `	// we are aiming for a single row so we will use Query Row	
	{{range $i, $e := .Columns}}var nullable{{$e.GoName}} {{getNullableType $e.GoType}} 
	{{end -}}
	
	{{- if not .IsReturnVoid}}{{if .IsReturnUserDefined}}{{$pointerSymbol := ""}}returnVal = new({{.ReturnGoType}}){{else}}{{$pointerSymbol := "&"}}{{end}}{{end}}
	
	err = currentDbHandle.QueryRow(JoinStringParts(queryParts,""), {{range $i, $e := .Parameters}}param{{.GoFriendlyName}}{{if ne (plus1 $i) $paramCount}},{{end}} {{end}})` +
	`{{if not .IsReturnVoid}}` +
	`.Scan({{if .IsReturnUserDefined}}{{$colCount := len .Columns}}` +
	`{{range $i, $e := .Columns}}&nullable{{$e.GoName}}{{if ne (plus1 $i) $colCount}},{{end}}{{end}}` +
	`{{else}}` +
	`&returnVal` +
	`{{end}})` +
	`{{else}}.Scan(){{end}}		
			
    switch {
    case err == ErrNoRows:
            // no such row found, return nil and nil
			err = nil
			return
    case err != nil:
            return
    default:
			{{if .IsReturnVoid}}
			// this function returns void, do not attempt doing anything else
			return
			{{else}}//BEGIN: non void operations	
										
			{{if not .IsReturnUserDefined}} //todo 
			{{else}}
			// BEGIN: assign any nullable values to the nullable fields inside the struct appropriately
			var isNullRecord bool = true
			{{range $i, $e := .Columns}}
				{{if .Nullable}}returnVal.Set{{.GoName}}(nullable{{$e.GoName}}.GetValue(), nullable{{$e.GoName}}.Valid)
				if nullable{{$e.GoName}}.Valid { isNullRecord = false }
				{{else}}returnVal.Set{{.GoName}}(nullable{{$e.GoName}}.GetValue())
				if nullable{{$e.GoName}}.Valid { isNullRecord = false }{{end}}
			{{end}}
			
			if isNullRecord == true {
				{{if .IsReturnUserDefined}}returnVal = nil
				{{else}}{{if not .IsReturnVoid}}isDbNull = true{{end}}
				{{end}}
			}
			
			// END: assign any nullable values to the nullable fields inside the struct appropriately	
			{{end}}
			
			return
			//END: non void operations{{end}}
    }	
`

const FUNCTION_TEMPLATE = `{{$paramCount := len .Parameters}}
{{$functionName := .GoFriendlyName}}
// Wrapper over the function named {{.DbName}}
{{if not .IsReturnASet}}{{if not .IsReturnUserDefined}}// For pure Go return types, a true isDbNull return parameter indicates that 
// the actual value returned from the database was nil, not the default value of the Go type{{end}}{{end}}
func (utilRef *tFunctionUtils) {{$functionName}}(` +
	`{{range $i, $e := .Parameters}}param{{.GoFriendlyName}} {{.GoType}}{{if ne (plus1 $i) $paramCount}},{{end}} {{end}})` +
	` ({{if not .IsReturnVoid}}returnVal {{if .IsReturnASet}}[]{{else}}{{if .IsReturnUserDefined}}*{{end}}{{end}}{{.ReturnGoType}},{{end}} err error{{if not .IsReturnVoid}}{{if not .IsReturnASet}}{{if not .IsReturnUserDefined}}, isDbNull bool{{end}}{{end}}{{end}}) {
						
	var errorPrefix = "tFunctionUtils.{{$functionName}}() ERROR: "
	
	currentDbHandle := GetDb()
	if currentDbHandle == nil {
		err = NewModelsErrorLocal(errorPrefix, "the database handle is nil")
		return
	}	
	
	// define the exec query
	var queryParts []string
	
	queryParts = append(queryParts, "SELECT * FROM ")
	queryParts = append(queryParts, "{{.DbName}}")
	//queryParts = append(queryParts, "( {{range $i, $e := .Parameters}}{{$e.DbName}} := ${{(plus1 $i)}}{{if ne (plus1 $i) $paramCount}},{{end}}{{end}} )")
	queryParts = append(queryParts, "( {{range $i, $e := .Parameters}}${{(plus1 $i)}}{{if ne (plus1 $i) $paramCount}},{{end}}{{end}} )")	
{{if .IsReturnASet }}
` + COMMON_CODE_FUNCTION_QUERY + `{{else}}
` + COMMON_CODE_FUNCTION_QUERYROW + `{{end}}	
}

`
