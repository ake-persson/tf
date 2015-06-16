# Range
{{range .Inp.Apples}}Apple: {{.}}
{{end}}

# Range 2
{{range .Inp.Apples}}Apple: {{.}}{{print ",\n"}}{{end}}

# Range 3
{{range $i, $e := .Inp.Apples}}Apple: {{$e}}{{if last $i $.Apples | not}}{{print ",\n"}}{{end}}{{end}}

# Join single-line
{{join ", " .Inp.Apples}}

# Join multi-line
{{join ",\n" .Inp.Apples}}
