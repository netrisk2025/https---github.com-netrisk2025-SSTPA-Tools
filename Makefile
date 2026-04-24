SHELL := /bin/bash

.PHONY: bootstrap backend-run backend-test reference-run reference-test frontend-dev frontend-build frontend-lint frontend-typecheck frontend-test compose-config copyright-check copyright-apply sbom-generate sbom-check verify

bootstrap:
	npm install

backend-run:
	cd backend && go run ./cmd/api

backend-test:
	cd backend && go test ./...

reference-run:
	cd tools/reference-pipeline && go run ./cmd/refstage --print-layout

reference-test:
	cd tools/reference-pipeline && go test ./...

frontend-dev:
	cd apps/desktop-shell && npm run dev

frontend-build:
	npm run workspaces:build

frontend-lint:
	npm run workspaces:lint

frontend-typecheck:
	npm run workspaces:typecheck

frontend-test:
	npm run workspaces:test

compose-config:
	docker compose -f infra/docker/compose.yaml config

copyright-check:
	cd tools/devtools/copyright && go run ./cmd/apply --check

copyright-apply:
	cd tools/devtools/copyright && go run ./cmd/apply --apply

sbom-generate:
	node tools/devtools/sbom/generate-sbom.mjs

sbom-check:
	node tools/devtools/sbom/generate-sbom.mjs --check

verify: copyright-check sbom-check backend-test reference-test frontend-lint frontend-typecheck frontend-test compose-config
