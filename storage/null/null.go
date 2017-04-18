package null

type NullStorage struct {
}

func (s NullStorage) GetCount() int {
	return 0
}
func (s NullStorage) SetCount(i int) {
	return
}

func NewNullStorage() (*NullStorage, error) {
	return &NullStorage{}, nil
}
