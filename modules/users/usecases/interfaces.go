package usecases

type HashPassword interface {
	Hash(passwordPlain string) (string, error)
	Verify(passwordPlain, passwordHashed string) (bool, error)
}
