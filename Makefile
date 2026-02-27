.PHONY: dev-garden dev-shrine build-garden build-shrine deploy-garden

dev-garden:
	cd garden && npm run dev

dev-shrine:
	cd shrine && go run .

build-garden:
	cd garden && npm run build

build-shrine:
	cd shrine && make build

deploy-garden:
	./scripts/deploy.sh