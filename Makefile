run:
	@templ generate
	@npx tailwind -i src/styles/tailwind.css -o public/tailwind.css
	@go run cmd/main.go
