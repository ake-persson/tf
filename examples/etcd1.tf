+{{ repeat 47 "-" }}+{{ repeat 22 "-" }}+
{{ printf "| %-45v | %-20v |" "host" "hwaddr" }}
+{{ repeat 47 "-" }}+{{ repeat 22 "-" }}+
{{ range $k, $e := .Etcd }}{{ printf "| %-45v | %-20v |\n" $k $e.interface.default.hwaddr }}{{ end }}+{{ repeat 47 "-" }}+{{ repeat 22 "-" }}+
