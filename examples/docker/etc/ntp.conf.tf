{{range .TimeServers}}server {{.}} minpoll 3 maxpoll 7
{{end}}

driftfile /var/lib/ntp/drift

enable stats
statsdir /var/log/ntpstats/
statistics loopstats peerstats
