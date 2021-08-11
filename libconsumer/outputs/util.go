package outputs

func Fail(err error) (Group, error) { return Group{}, err }

func Success(batchSize, retry int, clients ...Client) (Group, error) {
	return Group{
		Clients:   clients,
		BatchSize: batchSize,
		Retry:     retry,
	}, nil
}
