package logger

func (l Zap) Close() error {
	_ = l.Logger.Sync()

	if l.File != nil {
		err := l.File.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
