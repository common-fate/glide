devcli:
	go build -o bin/approvals cmd/devcli/main.go
	mv ./bin/approvals /usr/local/bin/

gdeploy:
	go build -o bin/gdeploy cmd/gdeploy/main.go
	mv ./bin/gdeploy /usr/local/bin/

generate:
	go generate ./...
	cd web && pnpm clienttypegen 
	pnpm prettier  --write **/openapi.yml
	