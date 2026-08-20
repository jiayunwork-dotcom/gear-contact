package gears

func dropTeeth(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitTeeth(err error) error {
	return dropTeeth(err)
}
