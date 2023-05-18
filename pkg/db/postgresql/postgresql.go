package postgresql

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mokan-r/golook/pkg/models"
)

type PostgreSQL struct {
	DB *pgxpool.Pool
}

func (p *PostgreSQL) Insert(commands models.Commands) error {
	stmt := `INSERT INTO commands (
                      path,
                      changed_file,
                      executed_command,
                      exit_code,
                      started_at,
                      finished_at
                      ) VALUES ($1, $2, $3, $4, $5, $6)`

	conn, err := p.DB.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		context.Background(),
		stmt,
		commands.Path,
		commands.ChangedFile,
		commands.ExecutedCommand,
		commands.ExitCode,
		commands.StartedAt,
		commands.FinishedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
