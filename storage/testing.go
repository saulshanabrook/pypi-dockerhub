package storage

type TestingClient int64

func NewTestingClient(initialTime int64) *TestingClient {
	tc := TestingClient(initialTime)
	return &tc
}

func (tc *TestingClient) SetTime(time int64) error {
	*tc = TestingClient(time)
	return nil
}

func (tc *TestingClient) GetTime() (int64, error) {
	return int64(*tc), nil
}
