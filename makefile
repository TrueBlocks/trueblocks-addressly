all:
	@wails dev 2>&1 | grep -v "AssetHandler"

lint:
	@staticcheck ./...
	@cd frontend && yarn run eslint . --ext .ts,.tsx

