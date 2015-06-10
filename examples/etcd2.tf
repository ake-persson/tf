{{ range $hk, $he := .Etcd }}
{{ $hk }}:
{{ range $k, $e := $he }}{{ if ismap $e | not }}  |- {{ printf "%s = %s\n" $k $e}}{{ end }}{{ end }}{{ end}}
