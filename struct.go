package main

// Corresponds to Trillian Log LeafValue
type Input struct {
	inputHash string
}

func newInput(inputHash string) *Input {
	return &Input{
		inputHash: inputHash,
	}
}

func (t *Input) Marshal() ([]byte, error) {
	return []byte(t.inputHash), nil
}

// Corresponds to Trillian Log ExtraData
type Output struct {
	name string
}

func newOutput(name string) *Output {
	return &Output{
		name: name,
	}
}

func (e *Output) Marshal() ([]byte, error) {
	return []byte(e.name), nil
}
