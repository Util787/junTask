## rest api for junior test task

## Main points of api:
Receiving all users should be with filters and pagination.

When creating a user 3 api requests from server must collect info about age, gender and nationality. Used apis:

https://api.agify.io/

https://api.genderize.io/

https://api.nationalize.io/

Updating parameters might be provided partically without changing the record entierly, but only parts that were provided.

## Installation
1.Clone repo

``` 
git clone https://github.com/Util787/junTask
```

2.Configure .env variables (ENV must be 'prod' or 'dev' or 'local')

Example:

```
ENV=local
SERVER_PORT=8000
DB_HOST=localhost
DB_PORT=5436
DB_USERNAME=postgres
DB_PASSWORD=1111
DB_NAME=postgres
SSLMODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=2222
REDIS_DB=0
```

3.Use migrate for your postgres db

```
-- migrate -path sql/schema -database "DB_URL?sslmode=disable" up
```

Example:

```
-- migrate -path sql/schema -database "postgres://postgres:1111@localhost:5436/postgres?sslmode=disable" up
```

4.Run by using next command in main directory:

```
go run cmd/main.go
```

5.All endpoints are described in API documentation (/swagger/index.html)
