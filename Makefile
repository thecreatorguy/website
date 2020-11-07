initmodules:
	git submodule init
	git submodule update
	cp ./modules/jumpybird/jumpybird.js ./assets/js/jumpybird.js

build:
	docker build -t itstimjohnson-website .

run: build
	docker-compose up -d

initdb:
	PGPASSWORD=verysecretpassword psql -h localhost -U postgres -f scripts/init.sql -a