package gears

func dropModule(err error) error {
	return err
}

func commitModule(err error) error {
	if err == nil {
		return nil
	}
	return dropModule(err)
}
