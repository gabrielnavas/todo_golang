dev:
	export DATABASE_USER=postgres && \
	export DATABASE_HOST=127.0.0.1 && \
	export DATABASE_PORT=5432 &&  \
	export DATABASE_PASSWORD=postgres && \
	export DATABASE_DBNAME=database && \
	export DATABASE_SSLMODE=disable && \
	export TOKEN_AUTH_SECRET_KEY=iamsecretkeyiamsecretkeyiamsecretkeyiamsecretkeyiamsecretkeyiamsecretkeyiamsecretkey \
	gin go run main.go 