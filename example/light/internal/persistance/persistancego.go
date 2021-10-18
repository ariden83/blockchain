package persistance

type Persistence struct {
	lastHash []byte
}

func (p *Persistence) GetLastHash() ([]byte, error) {
	return []byte{}, nil
}
func (p *Persistence) Update(lastHash []byte, hashSerialize []byte) error {
	return nil
}
func (p *Persistence) LastHash() []byte {
	return []byte{}
}
func (p *Persistence) GetCurrentHashSerialize(hash []byte) ([]byte, error) {
	return []byte{}, nil
}
func (p *Persistence) Close() {
}
func (p *Persistence) DBExists() bool {
	return true
}
func (p *Persistence) SetLastHash(lastHash []byte) {
	p.lastHash = lastHash
}
