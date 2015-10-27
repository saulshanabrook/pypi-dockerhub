package storage

type Client interface {
	SetTime(int64) error
	GetTime() (int64, error)
}
