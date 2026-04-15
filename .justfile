set windows-shell := ["pwsh.exe", "-c"]
set shell := ["/bin/bash", "-uc"]

default_svc := "book"
default_port := "8081"

default:
  just --list

serve svc=default_svc port=default_port:
    go run cmd/server/{{svc}}/main.go --port {{port}}

check:
  buf lint
  buf build

gen: check
  buf generate