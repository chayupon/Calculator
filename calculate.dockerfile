FROM postgres:12.4-alpine

COPY db/init.sh /docker-entrypoint-initdb.d/00_init.sh
COPY db/*.sql /docker-entrypoint-initdb.d/