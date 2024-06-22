run:
	@templ generate
	@npx tailwind -i src/styles/styles.css -o public/styles.css
	@go build -o tmp/main cmd/main.go
