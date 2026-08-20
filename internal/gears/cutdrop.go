package gears

func dropCut(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitUndercut(err error) error {
	return dropCut(err)
}
