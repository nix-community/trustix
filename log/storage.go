package log

type logStorage struct {
	// TODO: Implement persistent storage
	hashes [][]*Leaf
}

func (s *logStorage) LevelSize(level int) int {
	return len(s.hashes[level])
}

func (s *logStorage) Get(level int, index int) *Leaf {
	return s.hashes[level][index]
}

func (s *logStorage) SliceStart(level int) [][]*Leaf {
	return s.hashes[level:]
}

func (s *logStorage) SliceEnd(level int) [][]*Leaf {
	return s.hashes[:level]
}

func (s *logStorage) Size() int {
	return len(s.hashes)
}

func (s *logStorage) Append(level int, leaf *Leaf) {
	if s.Size() == level {
		// TODO: Figure out level "growth"
		h := []*Leaf{}
		s.hashes = append(s.hashes, h)
	}

	hashes := s.hashes[level]
	hashes = append(hashes, leaf)
	s.hashes[level] = hashes
}
