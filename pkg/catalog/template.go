// markdown-template.go: the markdown template used to template the sysl module
package catalog

const ProjectTemplate = `
{{/* Automatically generated by https://github.com/anz-bank/sysl-catalog it is strongly recommended not to edit this file */}}
{{range $name, $link := .Links}} [{{$name}}]({{$link}}) {{end}} | [Chat with us]({{.ChatLink}}) | [New bug or feature request]({{.FeedbackLink}})
# {{Base .Title}}

| Package |
----|{{range $val := Packages .Module}}
[{{$val}}]({{$val}}/README.md)|{{end}}

## Integration Diagram
<img src="{{CreateIntegrationDiagram .Module .Title false}}">

## End Point Analysis Integration Diagram
<img src="{{CreateIntegrationDiagram .Module .Title true}}">

`

const MacroPackageProject = `
{{/* Automatically generated by https://github.com/anz-bank/sysl-catalog it is strongly recommended not to edit this file */}}
[Chat with us]({{.ChatLink}}) | [New bug or feature request]({{.FeedbackLink}})
# {{Base .Title}}

| Package |
----|{{range $val := MacroPackages .Module}}
[{{$val}}]({{$val}}/README.md)|{{end}}

## Integration Diagram
<img src="{{CreateIntegrationDiagram .Module .Title false}}">

## End Point Analysis Integration Diagram
<img src="{{CreateIntegrationDiagram .Module .Title true}}">

`

const NewPackageTemplate = `
{{/* Automatically generated by https://github.com/anz-bank/sysl-catalog it is strongly recommended not to edit this file */}}
[Back](../README.md) | [Chat with us]({{ChatLink}}) | [New bug or feature request]({{FeedbackLink}})
{{$packageName := ModulePackageName .}}

# {{$packageName}}

## Integration Diagram
![]({{CreateIntegrationDiagram . $packageName false}})
{{$Apps := .Apps}}

{{$databases := false}}
{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{if and (eq (hasPattern $app.Attrs "ignore") false) (eq (hasPattern $app.Attrs "db") true)}}
{{$databases = true}}
{{end}}{{end}}

{{if $databases}}
## Database Index
| Database Application Name  | Source Location |
----|----{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{if and (eq (hasPattern $app.Attrs "ignore") false) (eq (hasPattern $app.Attrs "db") true)}}
[{{$appName}}](#Database-{{$appName}}) | [{{SourcePath $app}}]({{SourcePath $app}})|  {{end}}{{end}}
{{end}}

## Application Index
| Application Name | Method | Source Location |
----|----|----{{$Apps := .Apps}}{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{if eq (hasPattern $app.Attrs "ignore") false}}{{$Endpoints := $app.Endpoints}}{{range $endpointName := SortedKeys $Endpoints}}{{$endpoint := index $Endpoints $endpointName}}{{if eq (hasPattern $endpoint.Attrs "ignore") false}}
{{$appName}} | [{{$endpoint.Name}}](#{{$appName}}-{{SanitiseOutputName $endpoint.Name}}) | [{{SourcePath $app}}]({{SourcePath $app}})|  {{end}}{{end}}{{end}}{{end}}

## Type Index
| Application Name | Type Name | Source Location |
----|----|----{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{$types := $app.Types}}{{if ne (hasPattern $app.Attrs "db") true}}{{range $typeName := SortedKeys $types}}{{$type := index $types $typeName}}
{{$appName}} | [{{$typeName}}](#{{$appName}}.{{$typeName}}) | [{{SourcePath $type}}]({{SourcePath $type}})|{{end}}{{end}}{{end}}


{{if $databases}}
# Databases
{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}
{{if hasPattern $app.GetAttrs "db"}}

<details>
<summary>Database {{$appName}}</summary>

{{Attribute $app "description"}}
![]({{GenerateDataModel $app}})
</details>
{{end}}{{end}}
{{end}}

# Applications
{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}
{{if eq (hasPattern $app.Attrs "ignore") false}}
{{if eq (hasPattern $app.Attrs "db") false}}
{{if ne (len $app.Endpoints) 0}}

## Application {{$appName}}

- {{Attribute $app "description"}}

{{ServiceMetadata $app}}

{{with CreateRedoc $app.SourceContext $appName}}
[View OpenAPI Specs in Redoc]({{CreateRedoc $app.SourceContext $appName}})
{{end}}

{{range $e := $app.Endpoints}}
{{if eq (hasPattern $e.Attrs "ignore") false}}


### {{$appName}} {{SanitiseOutputName $e.Name}}
{{Attribute $e "description"}}

<details>
<summary>Sequence Diagram</summary>

![]({{CreateSequenceDiagram $appName $e}})
</details>

<details>
<summary>Request types</summary>

#### Request types
{{if and (eq (len $e.Param) 0) (not $e.RestParams) }}
No Request types
{{end}}

{{range $param := $e.Param}}
{{Attribute $param.Type "description"}}

![]({{CreateParamDataModel $app $param}})
{{end}}

{{if $e.RestParams}}{{if $e.RestParams.UrlParam}}
{{range $param := $e.RestParams.UrlParam}}
{{$pathDataModel := (CreateParamDataModel $app $param)}}
{{if ne $pathDataModel ""}}
#### Path Parameter

![]({{$pathDataModel}})
{{end}}{{end}}{{end}}

{{if $e.RestParams.QueryParam}}
{{range $param := $e.RestParams.QueryParam}}
{{$queryDataModel := (CreateParamDataModel $app $param)}}
{{if ne $queryDataModel ""}}
#### Query Parameter

![]({{$queryDataModel}})
{{end}}{{end}}{{end}}{{end}}
</details>
<details>
<summary>Response types</summary>

#### Response types
{{$responses := false}}
{{range $s := $e.Stmt}}{{$diagram := CreateReturnDataModel  $appName $s $e}}{{if ne $diagram ""}}
{{$responses = true}}
{{$ret := (GetReturnType $e $s)}}{{if $ret }}
{{Attribute $ret "description"}}{{end}}

![]({{$diagram}})

{{end}}{{end}}
{{if eq $responses false}}
No Response Types

{{end}}
</details>

---

{{end}}{{end}}{{end}}{{end}}{{end}}{{end}}


# Types

{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{$types := $app.Types}}
{{if ne (hasPattern $app.Attrs "db") true}}
{{range $typeName := SortedKeys $types}}{{$type := index $types $typeName}}
<details>
<summary>{{$appName}}.{{$typeName}}</summary>

### {{$appName}}.{{$typeName}}
{{$typedesc := (Attribute $type "description")}}
- {{if ne  $typedesc ""}}{{$typedesc}}{{end}}

![]({{CreateTypeDiagram  $appName $typeName $type false}})

[Full Diagram]({{CreateTypeDiagram  $appName $typeName $type true}})

#### Fields

| Field name | Type | Description |
|----|----|----|{{$fieldMap := Fields $type}}{{range $fieldName := SortedKeys $fieldMap}}{{$field := index $fieldMap $fieldName}}
| {{$fieldName}} | {{FieldType $field}} | {{$desc := Attribute $field "description"}}{{if ne $desc $typedesc}}{{$desc}}{{end}}|{{end}}

</details>{{end}}{{end}}{{end}}

<div class="footer">

`
