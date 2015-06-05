# Range
{{range .Apples}}Apple: {{.}}
{{end}}

# Range 2
{{range .Apples}}Apple: {{.}}{{print ",\n"}}{{end}}

# Range 3
{{range $i, $e := .Apples}}Apple: {{$e}}{{if last $i $.Apples | not}}{{print ",\n"}}{{end}}{{end}}

# Join single-line
{{join .Apples ", "}}

# Join multi-line
{{join .Apples ",\n"}}
