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

frontend public
    bind *:8000

    acl is_api_domain hdr_beg(host) -i api.
    http-request deny unless is_api_domain
    {{range $serviceName, $servers := .}}{{if not (eq $serviceName "bootstrap")}}
    acl is_{{$serviceName}} path_beg /{{$serviceName}}
    use_backend {{$serviceName}} if is_{{$serviceName}}
    {{end}}{{end}}
    default_backend bootstrap

{{range $serviceName, $servers := .}}
backend {{$serviceName}}
    {{range .}}
    server {{.Name}} {{.Ip}}:{{.Port}} check
    {{end}}
{{end}}