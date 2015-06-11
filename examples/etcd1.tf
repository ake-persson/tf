{{ define "border" }}+{{ repeat 47 "-" }}+{{ repeat 27 "-" }}+{{ end }}
{{ template "border" }}
| {{ "host" | lalign 45 }} | {{ "hwaddr" | center 25 }} |
{{ template "border" }}
{{ range $k, $e := .Etcd }}| {{ $k | lalign 45}} | {{ $e.interface.default.hwaddr | default "-" | center 25}} |
{{ end }}{{ template "border" }}
