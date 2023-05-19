![logo](misc/images/golook.svg)

<!-- TOC -->
* [Description](#description)
* [Requirements](#requirements)
* [Project structure](#project-structure)
* [Database schema](#database-schema)
<!-- TOC -->

# Description

This tool monitors directories and executes commands when files in these directories changes.
It's also works with regexp that you provide to exclude or include files to be monitored.
You can also provide path to a log file, to store logs from commands executions.
All info stores in PostgreSQL database (see [Database schema](#database-schema) for more info)

# Requirements

- [go1.20.2](https://pkg.go.dev/github.com/SunJary/dl/go1.20.2)
- [docker v20.10.23](https://docs.docker.com/engine/install/)
- [docker-compose v2.15.1](https://docs.docker.com/compose/install/linux/)

# Usage

First, you need to write a config file that provides
directories to be monitored and commands to be executed in this format:
```
- path: /home/user/project1
  commands:
    - go build -o ./build/bin/app1 cmd/service/main.go
    - go run ./build/bin/app1
```
You can pass multiple paths and commands to be executed for them. For example:
```
- path: /home/user/project1
  commands:
    - go build -o ./build/bin/app1 cmd/service/main.go
    - go run ./build/bin/app1
- path: /home/user/project2
  commands:
    - go test ./...
```
When changing files with documentation or tests, you do not want to restart the application.
So it's possible to enable/disable files based on a regular expression, for example:
```
- path: /home/user/project2
  include_regexp:
    - .*.go$
    - .*.env$
      exclude_regexp:
    - .*._test.go$
```


It's also possible to specify a log file where logs from command execution will be sent, for example:

```
- path: /home/user/project2
  log_file: /tmp/log2.out
```

After config file is done we need to run postgresql container using docker-compose. For convenience, I've created Makefile. So just run:
```bash
make build
```
This will install all the dependencies and build the application.


Then you can use golook as follows:
```bash
./golook --config path_to_your_config.yml
```

# Project structure

```
golook/
├── cmd/
│   └── golook/
│       └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── monitor/
│   │   ├── monitor.go
│   │   └── monitor_test.go
│   ├── executor/
│   │   ├── executor.go
│   │   └── executor_test.go
│   └── logger/
│       └── logger.go
├── migrations/
│   └── init_schema.sql
├── pkg/
│   └── db/
│       ├── models/
│       │   └── commands.go
│       ├── postgresql/
│       │   └── postgresql.go
│       └── db.go
├── misc/
│   └── images/
│       └── golook.svg
├── docker-compose.yml
├── go.mod
├── Makefile
└── README.md
```

# Database schema

Table commands:

| Column Name      | Data Type | Description                                                        |
|------------------|-----------|--------------------------------------------------------------------|
| id               | serial    | Unique identifier for the monitored directory                      |
| path             | text      | Absolute path to the monitored directory                           |
| changed_file     | text      | Absolute path to the changed file that triggered command execution |
| executed_command | text      | Command that have been executed                                    |
| exit_code        | int       | Exit code of the command execution                                 |
| created_at       | timestamp | Timestamp of the start of the command execution                    |
| updated_at       | timestamp | Timestamp of the end of the command execution                      |
