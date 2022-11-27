devcli:
	go build -o bin/commonfate cmd/devcli/main.go
	mv ./bin/commonfate /usr/local/bin/

gdeploy:
	go build -o bin/gdeploy cmd/gdeploy/main.go
	mv ./bin/gdeploy /usr/local/bin/

generate:
	go generate ./...
	cd web && pnpm clienttypegen 
	pnpm prettier  --write **/openapi.yml
	