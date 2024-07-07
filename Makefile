run: esbuild
	@templ generate
	@npx tailwind -i src/styles/styles.css -o public/styles.css
	@go build -o tmp/main cmd/main.go

esbuild:
	@npx esbuild src/scripts/main.ts --outfile=public/bundle.js --bundle --minify --sourcemap --target=chrome58,firefox57,safari11,edge16

fmt:
	@gofmt -w .
	@templ fmt .

test:
	@go test ./...
