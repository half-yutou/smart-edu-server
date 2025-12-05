package class_member

type Repository interface {
	CountMembers(classID int64) (int, error)
}
