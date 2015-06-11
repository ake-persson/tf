{{ define "border" }}+{{ repeat 47 "-" }}+{{ repeat 22 "-" }}+{{ end }}
{{ template "border" }}
{{ printf "| %-45v | %-20v |" "host" "hwaddr" }}
{{ template "border" }}
{{ range $k, $e := .Etcd }}{{ $hw := $e.interface.default.hwaddr | default "-" }}{{ printf "| %-45v | %-20v |\n" $k $hw }}{{ end }}{{ template "border" }}
