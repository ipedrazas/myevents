# myevents

REST api to manage my events

Database is Postgres:

```
-> %  docker run --name events-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
```

Data is populated using a script:

```
-> %  docker run -it --rm --link events-postgres:postgres postgres psql -h postgres -U postgres
```

