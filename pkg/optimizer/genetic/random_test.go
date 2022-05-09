package genetic

import (
	"math/rand"

	"github.com/golang/mock/gomock"
	mocks "github.com/lothar1998/v2x-optimizer/test/mocks/optimizer/genetic"
)

var commonRandom = rand.New(rand.NewSource(0))

type generatorStub struct {
	*mocks.MockRandomGenerator
	previousIntnCall *gomock.Call
	previousPermCall *gomock.Call
}

func newGeneratorStub() *generatorStub {
	generator := mocks.NewMockRandomGenerator(gomock.NewController(nil))
	return &generatorStub{MockRandomGenerator: generator}
}

func (g *generatorStub) WithNextInt(argument, value int) *generatorStub {
	if g.previousIntnCall != nil {
		g.previousIntnCall = g.EXPECT().Intn(gomock.Eq(argument)).Return(value).After(g.previousIntnCall)
	} else {
		g.previousIntnCall = g.EXPECT().Intn(gomock.Eq(argument)).Return(value)
	}
	return g
}

func (g *generatorStub) WithNextPermutation(permutation []int) *generatorStub {
	if g.previousPermCall != nil {
		g.previousPermCall = g.EXPECT().Perm(gomock.Eq(len(permutation))).Return(permutation).After(g.previousPermCall)
	} else {
		g.previousPermCall = g.EXPECT().Perm(gomock.Eq(len(permutation))).Return(permutation)
	}
	return g
}
