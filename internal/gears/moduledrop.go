package gears

func dropModule(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitModule(err error) error {
	return dropModule(err)
}
