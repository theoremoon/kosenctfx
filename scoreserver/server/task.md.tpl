## {{ .Name }}

{{range .Tags}} `{{- . -}}` {{ end }}

{{ .Description }}

author: {{ .Author }}

{{ $len := len .Attachments }}
{{ if gt $len 0 }}
### Attachments
{{ range $idx, $a := .Attachments }}
- [{{- $a.Name -}}]({{- $a.URL -}})
{{ end }}
{{ end }}

