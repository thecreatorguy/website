.PHONY: build run runprod

initmodules:
	git submodule init
	git submodule update --recursive --remote
	cp ./modules/jumpybird/jumpybird.js ./assets/js/jumpybird.js
	cp ./modules/shakesearch/completeworks.txt ./data/completeworks_shakespeare.txt
	
build:
	docker build -t itstimjohnson-website -f ./build/package/website.Dockerfile .
	docker build -t itstimjohnson-httpsredirect -f ./build/package/httpsredirect.Dockerfile .

run: build
	docker-compose -p itstimjohnson-website \
		-f ./build/package/docker-compose.yml \
		-f ./build/package/docker-compose.local.yml \
		up -d

push: build
	docker tag itstimjohnson-website thecreatorguy/website
	docker push thecreatorguy/website
	docker tag itstimjohnson-httpsredirect thecreatorguy/httpsredirect
	docker push thecreatorguy/httpsredirect

logs:
	docker-compose -p itstimjohnson-website \
		-f ./build/package/docker-compose.yml \
		-f ./build/package/docker-compose.local.yml \
		logs

stopprod:
	docker-compose -p itstimjohnson-website \
		-f ./build/package/docker-compose.yml \
		-f ./build/package/docker-compose.prod.yml \
		down

pullprod: stopprod
	docker-compose -p itstimjohnson-website \
		-f ./build/package/docker-compose.yml \
		-f ./build/package/docker-compose.prod.yml \
		pull

runprod: pullprod
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