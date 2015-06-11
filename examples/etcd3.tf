{{ $pad := $pad + "    " }}


{{ range $k, $e := .Etcd }}
{{ $k }}:
{{ range $k1, $e1 := $e }}{{ if ismap $e1 | not }}{{ printf "    |- %s = %s\n" $k1 $e1}}{{ end }}{{ if ismap $e1 }}{{ printf "    |- %s\n" $k1 }}{{ range $k2, $e2 := $e1 }}{{ if ismap $e2 | not }}{{ printf "        |- %s = %s\n" $k2 $e2}}{{ end }}{{ end }}{{ end }}{{ end }}{{ end}}
