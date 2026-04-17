set windows-shell := ["pwsh.exe", "-c"]
set shell := ["/bin/bash", "-uc"]

default_svc := "book"

default:
  just --list

serve svc=default_svc:
    go run cmd/server/{{svc}}/main.go

check:
  buf lint
  buf build

gen: check
  buf generate