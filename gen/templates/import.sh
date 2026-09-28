{{- if .IsBulk -}}
terraform import secureaccess_{{snakeCase .Name}}.example "{{range .Attributes}}{{if .Reference}}<{{.TfName}}>,{{end}}{{end}}[<item1_name>,<item2_name>,...]"
{{else -}}
terraform import secureaccess_{{snakeCase .Name}}.example "{{$id := false}}{{range .Attributes}}{{if .Id}}{{$id = true}}<{{.TfName}}>{{end}}{{end}}{{if not $id}}{{range .Attributes}}{{if .Reference}}<{{.TfName}}>,{{end}}{{end}}<id>{{end}}"
{{- end }}