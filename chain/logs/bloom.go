package logs

type Config struct {
	Levels          uint64
	ElementPerIndex uint64
}

type GroupPosition struct {
	/// Bloom level.
	Level uint64
	/// Index of the group.
	Index uint64
}

/// Uniquely identifies bloom position including the position in the group.
type Position struct {
	/// Group position.
	Group *GroupPosition
	/// Number in group.
	Number uint64
}
