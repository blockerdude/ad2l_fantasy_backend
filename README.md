## Database
Install [golang-migrate](https://github.com/golang-migrate/migrate)
If on windows; install [scoop](https://scoop.sh/) to help install.

Create migrations files

    migrate create -ext sql -dir .\sql\migrations\ -seq {migration_name}

Run migrations files

    migrate -source file://D:/git/ad2l_fantasy_backend/sql/migrations -database $Env:DATABASE_URL -verbose up



## Setup

A secrets.json file is required and has the following form

    {
        "db_conn_string": "postgres://{user}:{password}@{host}/{db_name}?sslmode=disable"
    }
As an example the conn string can look like `"postgres://myUser:myPwd@localhost/ad2l_fantasy?sslmode=disable"`
