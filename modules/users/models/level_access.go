package models

import "errors"

var (
	ErrMinLevelAccess = errors.New("level access is wrong")
	ErrMaxLevelAccess = errors.New("level access is wrong")
)

type LevelAccess int16

const minLevelAccess = 1
const maxLevelAccess = 3
const (
	BasicLevelAccess   LevelAccess = 1
	ManagerLevelAccess LevelAccess = 2
	AdminLevelAccess   LevelAccess = 3
)

func (l *LevelAccess) Valid() error {
	if *l < minLevelAccess {
		return ErrMinLevelAccess
	}
	if *l > maxLevelAccess {
		return ErrMaxLevelAccess
	}
	return nil
}
