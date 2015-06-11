{{ define "border" }}+{{ repeat 45 "-" }}+{{ repeat 25 "-" }}+{{ end }}
{{ template "border" }}
|{{ "host" | center 45 }}|{{ "hwaddr" | center 25 }}|
{{ template "border" }}
{{ range $k, $e := .Etcd }}|{{ $k | center 45}}|{{ $e.interface.default.hwaddr | default "-" | center 25}}|
{{ end }}{{ template "border" }}
