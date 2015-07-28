{{ define "border" }}+{{ repeat 22 "-" }}+{{ repeat 12 "-" }}+{{ end }}
{{ template "border" }}
| {{ "host" | lalign 20 }} | {{ "serialno" | lalign 10 }} |
{{ template "border" }}
{{ range $k, $e := .Hosts }}| {{ $k | lalign 20 }} | {{ $e.serialno | lalign 10 }} |
{{ end }}{{ template "border" }}
