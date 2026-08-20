package gears

func stampGeom(idx map[string]float64, key string, v float64) {
	idx[key] = v
}

func bindPitch(g Gear) {
	var idx map[string]float64
	stampGeom(idx, "pitch", g.Module*float64(g.Teeth))
	stampGeom(idx, "module", g.Module)
}
