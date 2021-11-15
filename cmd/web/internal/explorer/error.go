package explorer

type Error struct {
	Status    int
	Error     error
	PageTitle string
}
