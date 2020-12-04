.PHONY: build run prodrun

initmodules:
	git submodule init
	git submodule update
	cp ./modules/jumpybird/jumpybird.js ./assets/js/jumpybird.js

build:
	docker build -t itstimjohnson-website -f ./build/package/Dockerfile .

run: build
	docker-compose -p itstimjohnson-website -f ./build/package/docker-compose.yml up -d

runprod: httpsredirect
	docker-compose -p itstimjohnson-website \
		-f ./build/package/docker-compose.yml \
		-f ./build/package/docker-compose.prod.yml \
		up -d

prodlogs:
	docker-compose -p itstimjohnson-website \
		-f ./build/package/docker-compose.yml \
		-f ./build/package/docker-compose.prod.yml \
		logs

initdb:
	PGPASSWORD=verysecretpassword psql -h localhost -U postgres -f scripts/init.sql -a

sass:
	sass resources/css:assets/css

watch:
	sass --watch resources/css:assets/css

buildredirect:
	go build -o . cmd/https-redirect/https-redirect.go

httpsredirect:
	./https-redirect &