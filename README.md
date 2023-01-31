# Setup local environment


## Prepare database

```sh
docker-compose up -d


docker exec -it postgres bash
psql -U postgres postgresDB
select * from pg_available_extensions;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

go run migrate/migrate.go
```
