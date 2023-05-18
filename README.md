![logo](misc/images/golook.svg)

# Description

This tool monitors directories and executes commands when files in these directories changes.
Directories to be monitored and commands to be executed must be provided in the config in this format:
```
- path: /home/user/project1
  commands:
    - go build -o ./build/bin/app1 cmd/service/main.go
    - go run ./build/bin/app1
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

# TODO

- [ ] Write proper tests
- [ ] Better logging
- [ ] Write proper requirements for running application

# Requirements

- Migrate tool
```bash
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
│   ├── 001_init_schema.up.sql
│   └── 001_init_schema.down.sql
├── pkg/
│   └── db/
│       ├── models/
│       │   └── commands.go
│       ├── postgresql/
│       │   └── postgresql.go
│       ├── db.go
│       └── db_test.go
├── docker-compose.yml
├── go.mod
└── README.md
```

# Database schema

Table monitored_directories:

| Column Name      | Data Type | Description                                                        |
|------------------|-----------|--------------------------------------------------------------------|
| id               | serial    | Unique identifier for the monitored directory                      |
| path             | text      | Absolute path to the monitored directory                           |
| changed_file     | text      | Absolute path to the changed file that triggered command execution |
| executed_command | text      | Command that have been executed                                    |
| exit_code        | int       | Exit code of the command execution                                 |
| created_at       | timestamp | Timestamp of the start of the command execution                    |
| updated_at       | timestamp | Timestamp of the end of the command execution                      |
