FROM migrate/migrate:v4.15.2

WORKDIR /

COPY migrations migrations

ENTRYPOINT \
	migrate \
	-path migrations \
	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
	"up"
