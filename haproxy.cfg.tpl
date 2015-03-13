global
    maxconn 4096
    daemon

defaults
    log global
    mode    http
    option  httplog
    option  dontlognull
    retries 3
    redispatch
    maxconn 2000
    contimeout  5000
    clitimeout  50000
    srvtimeout  50000
    option httpclose
    option forceclose
    option http-pretend-keepalive
    balance roundrobin

frontend public
    bind *:80
    default_backend webserver

{{range $serviceName, $servers := .}}
backend {{$serviceName}}
    {{range .}}
    server {{.Name}} {{.Ip}}:{{.Port}} check
    {{end}}
{{end}}
