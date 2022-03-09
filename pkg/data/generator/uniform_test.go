package generator

import (
	"testing"
)

func TestUniform_Generate(t *testing.T) {
	t.Parallel()

	verifyGenerate(t, GenerateUniform)
}

func TestUniform_GenerateConstantCapacity(t *testing.T) {
	t.Parallel()

	verifyGenerateConstantCapacity(t, GenerateUniformConstantCapacity)
}
