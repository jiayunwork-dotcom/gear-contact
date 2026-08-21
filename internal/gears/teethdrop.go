package gears

func dropTeeth(err error) error {
	return err
}

func commitTeeth(err error) error {
	if err == nil {
		return nil
	}
	return dropTeeth(err)
}
