package moveController

type Mover interface {
	MoveX(pos int) error
	MoveXAbs(pos int) error
	CurrPosition() (x, y int)
	MoveY(pos int) error
	MoveYAbs(pos int) error
}
