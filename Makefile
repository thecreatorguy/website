build:
	docker build -t itstimjohnson-website .

run: build
	docker-compose up -d

initdb:
	PGPASSWORD=verysecretpassword psql -h localhost -f scripts/init.sql -a