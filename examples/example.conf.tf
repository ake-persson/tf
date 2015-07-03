# My templated fictitious configuration file

[main]
region = {{ .File.region }}
country = {{ .File.country }}

{{ range .File.list }}{{ . }},{{ end }}

{{ if eq .File.region "amer" }}
[amer]
city = New York
{{ end }}
