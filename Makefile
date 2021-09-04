db:
	sqlite3 sqlite.db < ./init.sql

build:
	go build -o=go-pass && ./go-pass start

docker-build:
	docker build \
		--rm \
		-f ./Dockerfile \
		-t serbanblebea/go-pass:0.1 \
		.

docker-run:
	docker run \
		-p 8081:8081 \
		--env-file ./.env \
		-d \
		--name go-pass \
		serbanblebea/go-pass:0.1

docker: docker-build docker-run

package:
	rm -r ./chrome-extension/popup/node_modules && \
	zip -r ./extension.zip ./chrome-extension

rm-package:
	rm ./extension.zip

unit-test:
	go test -v ./server/handler