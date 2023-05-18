![logo](misc/images/golook.svg)

# Description

Pass

# TODO

- [ ] Write proper tests
- [ ] Better logging

# Requirements

- Migrate tool
```bash
```

# Task

Problem

When writing code, there is often a need to perform certain actions immediately upon saving files. Depending on the situation, this could include:

- Building and running an application
- Running tests
- Building an application and deploying it to a server
- Running linters

Task

Implement a console application that allows you to:

- Monitor changes in various directories
- Execute an arbitrary set of console commands

Basic Requirements:

- It should be possible to specify a single directory to monitor.
- The application should be configurable via a configuration file, for example:

```
- path: /home/user/project1
  commands:
    - go build -o ./build/bin/app1 cmd/service/main.go
    - go run ./build/bin/app1
```
    
When changes are detected in /home/user/project1, the application should build and run the specified commands.
If one of the commands fails, the subsequent commands should not be executed.
History of file changes and command executions should be stored in a database.
When the application is stopped, all commands should be stopped as well, and some text (e.g. "finished") should be displayed.

Additional Features

It should be possible to monitor multiple directories, for example:

```
- path: /home/user/project1
  commands:
    - go build -o ./build/bin/app1 cmd/service/main.go
    - go run ./build/bin/app1
- path: /home/user/project2
  commands:
    - go test ./...
```

When changes are detected in /home/user/project1, the application should build and run the specified commands. When changes are detected in /home/user/project2, tests should be run.

It should be possible to enable/disable files based on a regular expression, for example:
```
- path: /home/user/project2
  include_regexp:
    - .*.go$
    - .*.env$
      exclude_regexp:
    - .*._test.go$
```
It should be possible to specify a log file where logs from command execution will be sent, for example:

```
- path: /home/user/project2
  log_file: /tmp/log2.out
```

Code Requirements

- Development language: Go
- Any libraries can be used.
- Relational database: PostgreSQL
- The code should be hosted on GitHub with a README file containing instructions for running the application and examples. It should be possible to follow the instructions and have the code working.
- If there are any questions about the requirements, the candidate can make decisions on their own. It is desirable to reflect in the README what questions were asked and what decisions were made.


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
│   ├── 001_init_schema.up.sql
│   └── 001_init_schema.down.sql
├── pkg/
│   └── db/
│       ├── db.go
│       └── db_test.go
├── docker-compose.yml
├── go.mod
└── README.md
```

# Database schema

Table monitored_directories:

| Column Name | Data Type | Description                                             |
|-------------|-----------|---------------------------------------------------------|
| id          | serial    | Unique identifier for the monitored directory           |
| path        | text      | Absolute path to the monitored directory                |
| log_file    | text      | Absolute path to the log file                           |
| created_at  | timestamp | Timestamp of the creation of the monitored directory    |
| updated_at  | timestamp | Timestamp of the last update of the monitored directory |
