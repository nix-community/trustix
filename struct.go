package main

// Input - Corresponds to Trillian Log LeafValue
type Input struct {
	inputHash  string
	outputHash string
}

// Create a new Input instance from inputHash string
func newInput(inputHash string, outputHash string) *Input {
	return &Input{
		inputHash:  inputHash,
		outputHash: outputHash,
	}
}

// Marshal - A convenience method
// While we probably don't want any extra fields in the input
// this method exists in case we want to extend the data type
// at some point later on
func (t *Input) Marshal() ([]byte, error) {

	return []byte(t.inputHash), nil
}
