// markdown-template.go: the markdown template used to template the sysl module
package catalog

const ProjectTemplateMermaid = `
<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>

{{/* Automatically generated by https://github.com/anz-bank/sysl-catalog it is strongly recommended not to edit this file */}}
{{range $name, $link := .Links}} [{{$name}}]({{$link}}) | {{end}} 
# {{Base .Title}}

| Package |
----|{{range $val := Packages .Module}}
[{{$val}}]({{$val}}/README.md)|{{end}}

## Integration Diagram
<pre class="mermaid">
{{IntegrationMermaid .Module .Title false}}
</pre>

## End Point Analysis Integration Diagram
<pre class="mermaid">
{{IntegrationMermaid .Module .Title true}}
</pre>

`

const MacroPackageProjectMermaid = `
<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>

{{/* Automatically generated by https://github.com/anz-bank/sysl-catalog it is strongly recommended not to edit this file */}}
# {{Base .Title}}

| Package |
----|{{if .Module}}{{range $val := MacroPackages .Module}}
[{{$val}}]({{$val}}/README.md)|{{end}}{{end}}

## Integration Diagram

<pre class="mermaid">
{{if .Module}}{{IntegrationMermaid .Module .Title false}}{{end}}
</pre>

## End Point Analysis Integration Diagram
<pre class="mermaid">
{{if .Module}}{{IntegrationMermaid .Module .Title true}}{{end}}
</pre>

`

const NewPackageTemplateMermaid = `
<script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>


{{/* Automatically generated by https://github.com/anz-bank/sysl-catalog it is strongly recommended not to edit this file */}}
[Back](../README.md)
{{$packageName := ModulePackageName .}}

# {{$packageName}}

## Integration Diagram
<pre class="mermaid">{{IntegrationMermaid . $packageName false}}</pre>

{{$Apps := .Apps}}

{{$databases := false}}
{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{if and (eq (hasPattern $app.Attrs "ignore") false) (eq (hasPattern $app.Attrs "db") true)}}
{{$databases = true}}
{{end}}{{end}}

{{if $databases}}
## Database Index
| Database Application Name  | Source Location |
----|----{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{if and (eq (hasPattern $app.Attrs "ignore") false) (eq (hasPattern $app.Attrs "db") true)}}
[{{SanitiseOutputName $appName}}](#Database-{{$appName}}) | [{{SourcePath $app}}]({{SourcePath $app}})|  {{end}}{{end}}
{{end}}

## Application Index
{{$anyApps := false}}

{{$Apps := .Apps}}{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{if eq (hasPattern $app.Attrs "ignore") false}}{{$Endpoints := $app.Endpoints}}{{range $endpointName := SortedKeys $Endpoints}}{{$endpoint := index $Endpoints $endpointName}}{{if eq (hasPattern $endpoint.Attrs "ignore") false}}{{if not $anyApps}}| Application Name | Method | Source Location |
|----|----|----|{{$anyApps = true}}{{end}}
| {{$appName}} | [{{$endpoint.Name}}](#{{SanitiseOutputName $appName}}-{{SanitiseOutputName $endpoint.Name}}) | [{{SourcePath $app}}]({{SourcePath $app}})|  {{end}}{{end}}{{end}}{{end}}

{{if not $anyApps}}
<span style="color:grey">No Applications Defined</span>
{{end}}


## Type Index
{{$anyTypes := false}}

{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{$types := $app.Types}}{{if ne (hasPattern $app.Attrs "db") true}}{{range $typeName := SortedKeys $types}}{{$type := index $types $typeName}}{{if not $anyTypes}}| Application Name | Type Name | Source Location |
|----|----|----|{{$anyTypes = true}}{{end}}
| {{$appName}} | [{{$typeName}}](#{{SanitiseOutputName $appName}}.{{SanitiseOutputName $typeName}}) | [{{SourcePath $type}}]({{SourcePath $type}})|{{end}}{{end}}{{end}}

{{if not $anyTypes}}
<span style="color:grey">No Types Defined</span>
{{end}}


{{if $databases}}
# Databases
{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}
{{if hasPattern $app.GetAttrs "db"}}

<a name=Database-{{SanitiseOutputName $appName}}></a><details>
<summary>Database {{$appName}}</summary>

{{Attribute $app "description"}}
<pre class="mermaid">
{{DataModelAppMermaid $app}}
</pre>

</details>
{{end}}{{end}}
{{end}}


{{if $anyApps}}
# Applications
{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}
{{if eq (hasPattern $app.Attrs "ignore") false}}
{{if eq (hasPattern $app.Attrs "db") false}}
{{if ne (len $app.Endpoints) 0}}

## Application {{$appName}}

{{$desc := Attribute $app "description"}}
{{if $desc}}
- {{$desc}}
{{end}}

{{ServiceMetadata $app}}

{{with CreateRedoc $app $appName}}
[View OpenAPI Specs in Redoc]({{CreateRedoc $app $appName}})
{{end}}

{{range $e := $app.Endpoints}}
{{if eq (hasPattern $e.Attrs "ignore") false}}


### <a name={{SanitiseOutputName $appName}}-{{SanitiseOutputName $e.Name}}></a>{{$appName}} {{$e.Name}}
{{Attribute $e "description"}}

<details>
<summary>Sequence Diagram</summary>

<pre class="mermaid">
{{SequenceMermaid $appName $e}}
</pre>
</details>

<details>
<summary>Request types</summary>

{{if and (not $e.Param) (not $e.RestParams) }}
<span style="color:grey">No Request types</span>
{{end}}
{{if not $e.Param}}{{if $e.RestParams }}{{if not $e.RestParams.UrlParam}}
<span style="color:grey">No Request types</span>
{{end}}{{end}}{{end}}

{{range $param := $e.Param}}
{{Attribute $param.Type "description"}}

<pre class="mermaid">
{{DataModelAliasMermaid $app $param}}
</pre>
{{end}}

{{if $e.RestParams}}{{if $e.RestParams.UrlParam}}
{{range $param := $e.RestParams.UrlParam}}
{{$pathDataModel := (DataModelAliasMermaid $app $param)}}
{{if ne $pathDataModel ""}}
#### Path Parameter

<pre class="mermaid">
{{$pathDataModel}}
</pre>
{{end}}{{end}}{{end}}

{{if $e.RestParams.QueryParam}}
{{range $param := $e.RestParams.QueryParam}}
{{$queryDataModel := (DataModelAliasMermaid $app $param)}}
{{if ne $queryDataModel ""}}
#### Query Parameter

<pre class="mermaid">
{{$queryDataModel}}
</pre>
{{end}}{{end}}{{end}}{{end}}
</details>

<details>
<summary>Response types</summary>

{{$responses := false}}
{{range $s := $e.Stmt}}{{$diagram := DataModelReturnMermaid  $appName $s $e}}{{if ne $diagram ""}}
{{$responses = true}}
{{$ret := (GetReturnType $e $s)}}{{if $ret }}
{{Attribute $ret "description"}}{{end}}

<pre class="mermaid">
{{$diagram}}
</pre>

{{end}}{{end}}

{{if not $responses}}
<span style="color:grey">No Response Types</span>
{{end}}
</details>
{{end}}

---

{{end}}{{end}}{{end}}{{end}}{{end}}{{end}}


{{if $anyTypes}}
# Types


{{range $appName := SortedKeys .Apps}}{{$app := index $Apps $appName}}{{$types := $app.Types}}
{{if ne (hasPattern $app.Attrs "db") true}}


{{range $typeName := SortedKeys $types}}{{$type := index $types $typeName}}
<a name={{SanitiseOutputName $appName}}.{{SanitiseOutputName $typeName}}></a><details>
<summary>{{$appName}}.{{$typeName}}</summary>

### {{$appName}}.{{$typeName}}
{{$typedesc := (Attribute $type "description")}}
{{if ne $typedesc ""}}- {{$typedesc}}{{end}}

<pre class="mermaid">
{{DataModelMermaid $appName $typeName}}
</pre>

[Full Diagram]({{DataModelMermaid $appName $typeName}})

{{if Fields $type}}
#### Fields
{{$fieldHeader := false}}
{{$fieldMap := Fields $type}}{{range $fieldName := SortedKeys $fieldMap}}{{$field := index $fieldMap $fieldName}}{{if not $fieldHeader}}| Field name | Type | Description |
|----|----|----|{{$fieldHeader = true}}{{end}}
| {{$fieldName}} | {{FieldType $field}} | {{$desc := Attribute $field "description"}}{{if ne $desc $typedesc}}{{$desc}}{{end}}|{{end}}
{{end}}

</details>{{end}}{{end}}{{end}}
{{end}}

<pre class="footer">

`
