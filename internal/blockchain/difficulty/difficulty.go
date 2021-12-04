package difficulty

type Difficulty int

var (
	Current      Difficulty = 1
	nbTryWaiting uint       = 500
)

func (d *Difficulty) Update(nbTries uint) {
	p := d.pourcent(nbTries)
	if p < 90 {
		*d = *d + 1
		return
	} else if p > 110 && *d > 2 {
		*d = *d - 1
		return
	}
}

func (Difficulty) pourcent(i uint) uint {
	return i * 100 / nbTryWaiting
}

func (d *Difficulty) Int() int {
	return int(*d)
}

func (d *Difficulty) Save(difficulty int) {
	*d = Difficulty(difficulty)
}
