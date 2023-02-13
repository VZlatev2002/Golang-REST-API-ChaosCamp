# Bill Splitter REST-API

This small project in part of the instructor-led remote course [Web programming with Go](https://www.chaos.com/chaos-camp/golang-2022), organized by the Chaos Group, a software company in Europe.

This project has only back-end as of now and the instructions on how to build it on a local machine are described hereafter.

## Building

### Prerequisite

The backend is written in Go and uses go modules (> go1.13).

### Setup

It uses MySQL as database and you'll need to create a new database using the schema in `./data/setup.sql`

The parameters listed below should be set as env vars:

`DB_USER` (default `mysql`)

`DB_PASSWORD` (default `password`)

`DB_NAME` (default `v1_bsplitter`)

`DB_PORT` (default `5432`)

`ServerAddr ` (default `:8080`)

To create the database, open MySQL CLI, log in with your credentials, run `create DB_NAME`, then `use DB_NAME`, and then `source /path/to/data/setup.sql`

To build the project run `go build -o .` and to run `/.bill-splitter`

## Example database scheme

![DB Scheme](/BillSplitter Golang.png)

## License

[MIT](https://choosealicense.com/licenses/mit/)
