package decoder

// addPadding adds 0-padding by creating a new slice which has length
// of multiple aes.Blocksize and fill it with input slice.
func addPadding(input []byte, blocksize int) []byte {
	numBytes := int(len(input)/blocksize+1) * blocksize
	newInput := make([]byte, numBytes)
	copy(newInput, input)

	return newInput
}
