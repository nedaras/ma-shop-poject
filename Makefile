run:
	@templ generate
	@npx tailwind -i src/styles/styles.css -o public/styles.css
	@go run cmd/main.go
