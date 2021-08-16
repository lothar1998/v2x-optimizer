package cmd

import (
	"github.com/golang/mock/gomock"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/test/mocks"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"math/rand"
	"path"
	"strconv"
	"testing"
)

func Test_generateWith(t *testing.T) {
	t.Parallel()

	t.Run("should generate single file", func(t *testing.T) {
		t.Parallel()

		expectedN := 10
		expectedV := 5

		dir, err := ioutil.TempDir("", "v2x-optimizer-generate-with-*")
		assert.NoError(t, err)
		testPath := path.Join(dir, "data_file_"+strconv.Itoa(int(rand.Uint32())))

		encoder := mocks.NewMockEncoderDecoder(gomock.NewController(t))
		encoder.EXPECT().
			Encode(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(data *data.Data, w io.Writer) error {
					assert.Len(t, data.MRB, expectedN)
					assert.Len(t, data.R, expectedV)
					_, err := w.Write([]byte{})
					assert.NoError(t, err)
					return nil
				})

		runnable := generateWith(encoder)

		command := &cobra.Command{}
		setUpGenerateFlags(command)
		err = command.Flags().Set(nValue, strconv.Itoa(expectedN))
		assert.NoError(t, err)
		err = command.Flags().Set(vValue, strconv.Itoa(expectedV))
		assert.NoError(t, err)
		err = command.Flags().Set(timesValue, strconv.Itoa(1))
		assert.NoError(t, err)

		err = runnable(command, []string{testPath})
		assert.NoError(t, err)
		assert.FileExists(t, testPath)
	})

	t.Run("should generate multiple files", func(t *testing.T) {
		t.Parallel()

		expectedN := 12
		expectedV := 4

		dir, err := ioutil.TempDir("", "v2x-optimizer-generate-with-*")
		assert.NoError(t, err)
		testPath := path.Join(dir, "data_file_"+strconv.Itoa(int(rand.Uint32())))

		encoder := mocks.NewMockEncoderDecoder(gomock.NewController(t))
		encoder.EXPECT().
			Encode(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(data *data.Data, w io.Writer) error {
					assert.Len(t, data.MRB, expectedN)
					assert.Len(t, data.R, expectedV)
					_, err := w.Write([]byte{})
					assert.NoError(t, err)
					return nil
				}).
			Times(2)

		runnable := generateWith(encoder)

		command := &cobra.Command{}
		setUpGenerateFlags(command)
		err = command.Flags().Set(nValue, strconv.Itoa(expectedN))
		assert.NoError(t, err)
		err = command.Flags().Set(vValue, strconv.Itoa(expectedV))
		assert.NoError(t, err)
		err = command.Flags().Set(timesValue, strconv.Itoa(2))
		assert.NoError(t, err)

		err = runnable(command, []string{testPath + ".txt"})
		assert.NoError(t, err)
		assert.FileExists(t, testPath+"_0.txt")
		assert.FileExists(t, testPath+"_1.txt")
	})
}

func Test_generateDataFile(t *testing.T) {
	t.Parallel()

	t.Run("should generate single file with given filepath", func(t *testing.T) {
		t.Parallel()

		expectedContent := "this is expected content"
		expectedN := 10
		expectedV := 5

		dir, err := ioutil.TempDir("", "v2x-optimizer-generate-data-file-*")
		assert.NoError(t, err)

		testPath := path.Join(dir, "data_file_"+strconv.Itoa(int(rand.Uint32())))

		encoder := mocks.NewMockEncoderDecoder(gomock.NewController(t))
		encoder.EXPECT().
			Encode(gomock.Any(), gomock.Any()).
			DoAndReturn(
				func(data *data.Data, w io.Writer) error {
					assert.Len(t, data.MRB, expectedN)
					assert.Len(t, data.R, expectedV)
					_, err := w.Write([]byte(expectedContent))
					assert.NoError(t, err)
					return nil
				})

		err = generateDataFile(testPath, encoder, uint(expectedN), uint(expectedV))

		assert.NoError(t, err)
		assert.FileExists(t, testPath)

		readData, err := ioutil.ReadFile(testPath)
		assert.NoError(t, err)

		assert.Equal(t, expectedContent, string(readData))
	})
}

func Test_toMultipleFilesFilepath(t *testing.T) {
	t.Parallel()

	type args struct {
		path string
		i    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"should create filepath for file", args{"file", 2}, "file_2"},
		{"should create filepath for file with extension", args{"file.txt", 2}, "file_2.txt"},
		{"should create filepath for filepath", args{"/dir1/dir2/file", 2}, "/dir1/dir2/file_2"},
		{"should create filepath for filepath with extension", args{"/dir1/dir2/file.txt", 2},
			"/dir1/dir2/file_2.txt"},
		{"should create filepath for same filename and extension",
			args{"./third_party/cplex/data/data.dat", 2}, "./third_party/cplex/data/data_2.dat"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toMultipleFilesFilepath(tt.args.path, tt.args.i); got != tt.want {
				t.Errorf("toMultipleFilesFilepath() = %v, want %v", got, tt.want)
			}
		})
	}
}
