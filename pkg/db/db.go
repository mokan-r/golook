package db

import "github.com/mokan-r/golook/pkg/models"

type DB interface {
	Insert(commands models.Commands) error
}
