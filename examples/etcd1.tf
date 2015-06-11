{{ define "border" }}+{{ repeat 47 "-" }}+{{ repeat 27 "-" }}+{{ repeat 27 "-" }}+{{ repeat 42 "-" }}+{{ end }}
{{ template "border" }}
| {{ "host" | lalign 45 }} | {{ "hwaddr" | center 25 }} | {{ "ip" | center 25 }} | {{ "image" | lalign 40 }} |
{{ template "border" }}
{{ range $k, $e := .Etcd }}| {{ $k | lalign 45}} | {{ $e.interface.default.hwaddr | default "-" | center 25}} | {{ $e.interface.default.ip | default "-" | center 25 }} | {{ $e.image | default "-" | lalign 40 }} |
{{ end }}{{ template "border" }}
