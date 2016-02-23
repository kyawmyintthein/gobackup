# gobackup
Database backup command line tools for Golang
This libray is mainly based on [Cli](http://github.com/codegangsta/cli) and [Barkup](https://github.com/kyawmyintthein/barkup).

**Adapter:** `MySQL` `Postgres` `Mongodb`

**Target:** `S3`

This command line tool is that backup the database (mysql, postgres and mongodb) to amazon s3. It is congiurable by using yml file.

**Usage**
```go
go get "github.com/kyawmyintthein/gobackup"
cd go/src/kyawmyintthein/gobackup 
go install
gobackup export -f='path/to/gobackup.yml'

```

```gobackup.yml

adapter: mysql

database:
    name:     test
    user:     test
    password: test
    port:     1234
    host:     localhost

s3:
    bucket: backups
    region: us-east-1
    access_key: 122223dsfgadg
    secret: 122223dsfgadg
    path: data_backup/

```

Feel free to contribute this libray.


**Road Map**
***Target Storage Source***
Local file storage 

***Featured***
Cron job for auto backup


