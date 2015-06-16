# My templated fictitious configuration file

[main]
region = {{ .Inp.region }}
country = {{ .Inp.country }}

{{ range .Inp.list }}{{ . }},{{ end }}

{{ if eq .Inp.region "amer" }}
[amer]
city = New York
{{ end }}
