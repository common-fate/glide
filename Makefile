devcli:
	go build -o bin/commonfate cmd/devcli/main.go
	mv ./bin/commonfate /usr/local/bin/

gdeploy:
	go build -o bin/gdeploy cmd/gdeploy/main.go
	mv ./bin/gdeploy /usr/local/bin/

generate:
	go generate ./...
	cd web && pnpm clienttypegen 
	cd deploymentcli/web && pnpm clienttypegen
	pnpm prettier  --write **/openapi.yml
	pnpm prettier  --write ./deploymentcli.openapi.yml
	
clean:
	pnpm prettier  --write **/openapi.yml
	pnpm prettier  --write ./deploymentcli.openapi.yml
