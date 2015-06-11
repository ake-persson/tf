{{define "border"}}+{{ repeat 47 "-" }}+{{ repeat 22 "-" }}+{{end}}
{{template "border"}}
{{ printf "| %-45v | %-20v |" "host" "hwaddr" }}
{{template "border"}}
{{ range $k, $e := .Etcd }}{{ printf "| %-45v | %-20v |\n" $k $e.interface.default.hwaddr }}{{ end }}{{ template "border" }}
