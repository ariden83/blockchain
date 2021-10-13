package wallet

// Seed represents each 'item' in the blockchain
type Seed struct {
	Address   string
	Timestamp string
	PubKey    string
	PrivKey   string
	Mnemonic  string
}

type SeedNoPrivKey struct {
	Address   string
	Timestamp string
	PubKey    string
}

func (ws *Wallets) GetAllSeeds() []SeedNoPrivKey {
	var allSeeds []SeedNoPrivKey
	for _, j := range ws.Seeds {
		allSeeds = append(allSeeds, SeedNoPrivKey{
			Address:   j.Address,
			Timestamp: j.Timestamp,
			PubKey:    j.PubKey,
		})
	}
	return allSeeds
}
