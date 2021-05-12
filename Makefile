# Start all services, monitoring and the web UI. Does not start server clients
start:
	docker-compose up -d usersrv customersrv productsrv promotionsrv auditsrv
	docker-compose up -d grafana
	docker-compose up web
stop:
	docker-compose down

# Docker-compose sample commands
composeup:
	docker-compose up
composedown:
	docker-compose down

# BuildPacks

packbuildall:
	make packbuildsrv
	make packbuildweb
	make packbuildtestlclients

packbuildweb:
	APIURL='http://localhost:8080/' npm run build  --prefix web/sapper
	pack build cdconweb --path ./web/sapper  --builder gcr.io/buildpacks/builder:v1

packbuildsrv:
	pack build cdconusersrv --env BP_GO_TARGETS=./user/server
	pack build cdconpromotionsrv --env BP_GO_TARGETS=./promotion/server
	pack build cdconproductsrv --env BP_GO_TARGETS=./product/server
	pack build cdconcustomersrv --env BP_GO_TARGETS=./customer/server
	pack build cdconauditsrv --env BP_GO_TARGETS=./audit/server

packbuildtestlclients:
	pack build cdconusercli --env BP_GO_TARGETS=./user/client
	pack build cdconpromotioncli --env BP_GO_TARGETS=./promotion/client
	pack build cdconproductcli --env BP_GO_TARGETS=./product/client
	pack build cdconcustomercli --env BP_GO_TARGETS=./customer/client

# Running individual services examples

# Build and start service (SERVICE value should be lowercase)
buildsrvdev:
	pack build cdcon$(SERVICE)srv --env GOOGLE_BUILDABLE=./$$SERVICE/server --builder gcr.io/buildpacks/builder:v1
buildsrv:
	pack build cdcon$(SERVICE)srv --env BP_GO_TARGETS=./$$SERVICE/server
runsrv:
	docker-compose up $(SERVICE)srv

# Build and start testing client (SERVICE value should be lower case)
buildclidev:
	pack build cdcon$(SERVICE)cli --env GOOGLE_BUILDABLE=./$$SERVICE/client --builder gcr.io/buildpacks/builder:v1
buildcli:
	pack build cdcon$(SERVICE)cli --env BP_GO_TARGETS=./$(SERVICE)/client
runcli:
	docker-compose up $(SERVICE)cli

#DockerHub
hubpush:
	pack build bolbeck/cdcon$(SERVICE)srv --env BP_GO_TARGETS=./$$SERVICE/server
	docker push bolbeck/cdcon$(SERVICE)srv

hubpushweb:
	#echo "no implemented yet"
	APIURL='http://localhost:8080/' npm run build  --prefix web/sapper
	pack build bolbeck/cdconweb --path ./web/sapper  --builder gcr.io/buildpacks/builder:v1
	docker push bolbeck/cdconweb

# -------------------------------------------------------------------------------------

# Run service directly
runpromosrv:
	go run promotion/server/promotionServer.go
runpromocli:
	go run promotion/client/promotionClient.go

# -------------------------------------------------------------------------------------

# Web App
# Directly (dev)
runweb:
	npm run dev

# Docker
docrunweb:
	docker run --rm --name cdconwebcont -p 3000:8080 cdconweb

#Docker-compose
composeupweb:
	docker-compose up web

# -------------------------------------------------------------------------------------

# Compile proto files
genpromotionproto:
	protoc --proto_path=$$GOPATH/src:. --micro_out=source_relative:.. --go_out=. --go_opt=paths=source_relative promotion/proto/promotion.proto
genuserproto:
	protoc --proto_path=$$GOPATH/src:. --micro_out=source_relative:.. --go_out=. --go_opt=paths=source_relative user/proto/user.proto
gencustomerproto:
	protoc --proto_path=$$GOPATH/src:. --micro_out=source_relative:.. --go_out=. --go_opt=paths=source_relative customer/proto/customer.proto
genproductproto:
	protoc --proto_path=$$GOPATH/src:. --micro_out=source_relative:.. --go_out=. --go_opt=paths=source_relative product/proto/product.proto
genstandardFieldsproto:
	protoc --proto_path=$$GOPATH/src:. --micro_out=source_relative:.. --go_out=. --go_opt=paths=source_relative globalProtos/standardFields.proto

# -------------------------------------------------------------------------------------

# Call service through the micro gateway
promoviaapigateway:
	curl --location --request POST 'http://localhost:8080/promotion/promotionSrv/getPromotions' \
    --header 'Content-Type: application/json' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjp7ImlkIjoyMzQzNzI1MzkxMjkxNjE4MzA1LCJjb21wYW55IjoiRHVjayBJbmMuIn0sImV4cCI6MTU5NzMzNTMzNywiaWF0IjoxNTk3MjQ4OTM3LCJpc3MiOiJnb1RlbXAudXNlcnNydiJ9.QWAvvoXQHv_Cf48PTrjK9uRvrdEblNvFOxQWjNcX79U' \
    --data-raw '{"name":"Promo1", "customerId": "ducksrus"}'

# Call service using the micro gateway running behind the ingress in K8s
authviaapigateway:
	curl --location --request POST 'http://gotemp.tst/user/userSrv/auth' \
	--header 'Content-Type: application/json' \
	--data-raw '{"pwd":"1234","email":"duck@mymail.com"}'

# ----  Monitoring --------

# Check the service metrics when running on docker
getsrvmetrics:
	curl http://localhost:$$PORT/metrics