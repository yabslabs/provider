package storage

type Storage interface {
	Save(path, filename string, data []byte) error
}

func NewStorage() Storage {
	return &FileSystem{}
}
