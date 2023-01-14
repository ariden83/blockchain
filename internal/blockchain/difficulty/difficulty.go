package difficulty

// Difficulty represent a mining difficulty.
type Difficulty int

var (
	// Current represent the actual mining difficulty.
	Current      Difficulty = 1
	nbTryWaiting uint       = 500
)

// Update the difficulty according to the number of tests that were carried out before finding the last block.
func (d *Difficulty) Update(nbTries uint) {
	p := d.percent(nbTries)
	if p < 90 {
		*d = *d + 1
		return
	} else if p > 110 && *d > 2 {
		*d = *d - 1
		return
	}
}

// percent compares the number of successful tests to mine the block with the average reference wanted to unlock a block.
func (Difficulty) percent(i uint) uint {
	return i * 100 / nbTryWaiting
}

func (d *Difficulty) Int() int {
	return int(*d)
}

func (d *Difficulty) Save(difficulty int) {
	if difficulty <= 1 {
		return
	}
	*d = Difficulty(difficulty)
}
