.PHONY: build frontend geo backend

build: frontend geo backend

frontend:
	cd web && pnpm install && pnpm build

geo:
	chmod +x download_geolite.sh
	./download_geolite.sh

backend:
	go build -o main ./cmd/api

clean:
	rm -rf main public/dist GeoLite2-City
