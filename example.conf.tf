# My templated fictitious configuration file

[main]
region = {{.region}}
country = {{.country}}

{{if eq .region "amer"}}
[amer]
city = New York
{{ end }}
