# Range
{{$last := len .Apples}}{{range .Apples}}Apple: {{.}}
{{end}}

# Join single-line
{{join .Apples ", "}}

# Join multi-line
{{join .Apples ",\n"}}
