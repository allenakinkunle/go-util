package identifier

import (
	"github.com/gofrs/uuid"
	gonanoid "github.com/matoous/go-nanoid"
)

type IIdentifier interface {
	NewNanoID() string
	NewUUIDv5(name string) string
}

type identifier struct {
	characters string
	length     uint
	namespace  string
}

func NewIdentifier(characters string, length uint, namespace string) (*identifier, error) {
	return &identifier{
		characters: characters,
		length:     length,
		namespace:  namespace,
	}, nil
}

func (i identifier) NewNanoID() string {
	id, _ := gonanoid.Generate(i.characters, int(i.length))
	return id
}

func (i identifier) NewUUIDv5(name string) string {
	return uuid.NewV5(uuid.FromStringOrNil(i.namespace), name).String()
}
