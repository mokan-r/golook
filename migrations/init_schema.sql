CREATE TABLE IF NOT EXISTS commands (
                                        id SERIAL PRIMARY KEY,
                                        path TEXT NOT NULL,
                                        changed_file TEXT NOT NULL,
                                        executed_command TEXT NOT NULL,
                                        exit_code INT NOT NULL,
                                        started_at TIMESTAMP NOT NULL,
                                        finished_at TIMESTAMP NOT NULL
);
