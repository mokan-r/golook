CREATE TABLE IF NOT EXISTS commands (
                                        id SERIAL PRIMARY KEY,
                                        path TEXT NOT NULL,
                                        changed_file TEXT NOT NULL,
                                        executed_command TEXT NOT NULL,
                                        exit_code INT NOT NULL,
                                        started_at TIMESTAMP NOT NULL,
                                        finished_at TIMESTAMP NOT NULL
);
-- CREATE TABLE IF NOT EXISTS monitored_files (
--                                                id SERIAL PRIMARY KEY,
--                                                monitored_dir_id INTEGER NOT NULL,
--                                                path TEXT NOT NULL,
--                                                is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
--                                                include_regexp TEXT,
--                                                exclude_regexp TEXT,
--                                                created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--                                                updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--                                                FOREIGN KEY (monitored_dir_id) REFERENCES monitored_directories (id) ON DELETE CASCADE
-- );
-- CREATE TABLE IF NOT EXISTS commands (
--                                         id SERIAL PRIMARY KEY,
--                                         monitored_dir_id INTEGER NOT NULL,
--                                         command TEXT NOT NULL,
--                                         created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--                                         updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--                                         FOREIGN KEY (monitored_dir_id) REFERENCES monitored_directories (id) ON DELETE CASCADE
-- );
--
-- CREATE TABLE IF NOT EXISTS command_executions (
--                                                   id SERIAL PRIMARY KEY,
--                                                   command_id INTEGER NOT NULL,
--                                                   output TEXT,
--                                                   exit_code INTEGER,
--                                                   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
--                                                   FOREIGN KEY (command_id) REFERENCES commands (id) ON DELETE CASCADE
-- );
