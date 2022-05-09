package genoperator

type RandomGenerator interface {
	Intn(int) int
	Perm(int) []int
}
