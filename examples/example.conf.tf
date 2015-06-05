# My templated fictitious configuration file

[main]
region = {{.region}}
country = {{.country}}

{{range .list}}{{.}},{{end}}

{{if eq .region "amer"}}
[amer]
city = New York
{{ end }}
