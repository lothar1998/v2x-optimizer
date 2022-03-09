package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"path"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	encoderMock "github.com/lothar1998/v2x-optimizer/test/mocks/data"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_generateWith(t *testing.T) {
	t.Parallel()

	type args struct {
		times int
	}
	tests := []struct {
		name string
		args args
	}{
		{"should generate single file", args{1}},
		{"should generate multiple files", args{3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectedContent := "this is expected content"
			expectedN := 12
			expectedV := 4

			dir, err := ioutil.TempDir("", "v2x-optimizer-generate-with-*")
			assert.NoError(t, err)

			filename := "data_file_" + strconv.Itoa(int(rand.Uint32()))
			testPath := path.Join(dir, filename)

			encoder := encoderMock.NewMockEncoderDecoder(gomock.NewController(t))
			encoder.EXPECT().
				Encode(gomock.Any(), gomock.Any()).
				DoAndReturn(
					func(data *data.Data, w io.Writer) error {
						assert.Len(t, data.MRB, expectedN)
						assert.Len(t, data.R, expectedV)
						_, err := w.Write([]byte(expectedContent))
						assert.NoError(t, err)
						return nil
					}).
				Times(tt.args.times)

			runnable := generateWith(encoder)

			command := &cobra.Command{}
			setUpGenerateFlags(command)
			err = command.Flags().Set(bucketCountValue, strconv.Itoa(expectedN))
			assert.NoError(t, err)
			err = command.Flags().Set(itemCountValue, strconv.Itoa(expectedV))
			assert.NoError(t, err)
			err = command.Flags().Set(timesValue, strconv.Itoa(tt.args.times))
			assert.NoError(t, err)

			err = runnable(command, []string{testPath})
			assert.NoError(t, err)

			for i := 0; i < tt.args.times; i++ {
				expectedFilepath := filepath.Join(testPath, fmt.Sprintf(outputFilePattern, i))
				assert.FileExists(t, expectedFilepath)

				readData, err := ioutil.ReadFile(expectedFilepath)
				assert.NoError(t, err)
				assert.Equal(t, expectedContent, string(readData))
			}
		})
	}
}
