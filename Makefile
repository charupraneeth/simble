GOPATH=$(shell go env GOPATH)
AIR=$(GOPATH)/bin/air

.PHONY: build frontend geo backend dev server

build: frontend geo backend

frontend:
	cd web && pnpm install && pnpm build

geo:
	chmod +x download_geolite.sh
	./download_geolite.sh

backend:
	go build -o main ./cmd/api

dev:
	make -j2 server frontend-dev

server:
	$(AIR)

frontend-dev:
	cd web && pnpm dev

clean:
	rm -rf main public/dist GeoLite2-City tmp
