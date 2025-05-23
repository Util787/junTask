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
git clone github.com/Util787/junTask
```

2.Configure .env variables
Example:

```
SERVERPORT=8000
DBHOST=localhost
DBPORT=5436
DBUSERNAME=postgres
DBPASSWORD=1111
DBNAME=postgres
SSLMODE=disable
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