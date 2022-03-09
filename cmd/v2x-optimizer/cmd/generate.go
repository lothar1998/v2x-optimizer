package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/data/generator"
	"github.com/spf13/cobra"
)

const (
	itemCountValue          = "items"
	itemSizeValue           = "item-size"
	bucketCountValue        = "buckets"
	bucketSizeValue         = "bucket-size"
	constantBucketSizeValue = "constant-bucket-size"
	countValue              = "count"
	kindValue               = "kind"

	exponentialKind = "exponential"
	normalKind      = "normal"
	uniformKind     = "uniform"
	v2xKind         = "v2x"

	outputFilePattern = "data_%d.v2x"
)

type generateFunc func(itemCount, maxItemSize, bucketCount, bucketSize int) *data.Data

// GenerateCmd returns cobra.Command which is able to generate data in specified format.
// It should be registered in root command using AddCommand() method.
func GenerateCmd() *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate data in given format",
		Long:  "Allows generating data in given format",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	for formatName, encoderInfo := range formatsToEncodersInfo {
		generateToCmd := generateTo(formatName, encoderInfo)
		setUpGenerateFlags(generateToCmd)
		generateCmd.AddCommand(generateToCmd)
	}

	return generateCmd
}

func generateTo(formatName string, encoderInfo encoderInfo) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("%s {output_file}", formatName),
		Args:  cobra.ExactArgs(1),
		Short: fmt.Sprintf("Generate data in %s format", encoderInfo.FormatDisplayName),
		Long:  fmt.Sprintf("Allows for generating data in %s format", encoderInfo.FormatDisplayName),
		RunE:  generateWith(encoderInfo.Encoder),
	}
}

func generateWith(encoder data.EncoderDecoder) func(*cobra.Command, []string) error {
	return func(command *cobra.Command, args []string) error {
		output := args[0]

		itemCount, err := command.Flags().GetUint(itemCountValue)
		if err != nil {
			return err
		}

		itemSize, err := command.Flags().GetUint(itemSizeValue)
		if err != nil {
			return err
		}

		bucketCount, err := command.Flags().GetUint(bucketCountValue)
		if err != nil {
			return err
		}

		bucketSize, err := command.Flags().GetUint(bucketSizeValue)
		if err != nil {
			return err
		}

		kind, err := command.Flags().GetString(kindValue)
		if err != nil {
			return err
		}

		isBucketSizeConstant, err := command.Flags().GetBool(constantBucketSizeValue)
		if err != nil {
			return err
		}

		generate := toGeneratorFunc(kind, isBucketSizeConstant)
		if generate == nil {
			return errors.New("unknown kind")
		}

		count, err := command.Flags().GetUint(countValue)
		if err != nil {
			return err
		}

		for i := uint(0); i < count; i++ {
			err := generateDataFile(output, int(i), encoder, itemCount, itemSize, bucketCount, bucketSize, generate)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func generateDataFile(
	outputPath string,
	fileID int,
	encoder data.EncoderDecoder,
	itemCount, itemSize, bucketCount, bucketSize uint,
	generate generateFunc,
) error {
	err := os.MkdirAll(outputPath, 0775)
	if err != nil {
		return fmt.Errorf("%w: %s", errCannotCreatePath, outputPath)
	}

	outputPath, err = filepath.Abs(outputPath)
	if err != nil {
		return err
	}

	filePath := path.Join(outputPath, fmt.Sprintf(outputFilePattern, fileID))

	outputFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("%w: %s", errCannotOpenFile, outputPath)
	}
	defer outputFile.Close()

	generatedData := generate(int(itemCount), int(itemSize), int(bucketCount), int(bucketSize))

	err = encoder.Encode(generatedData, outputFile)
	if err != nil {
		return fmt.Errorf("%w: %s", errCannotEncodeData, err.Error())
	}

	return nil
}

func toGeneratorFunc(kind string, isBucketSizeConstant bool) generateFunc {
	switch {
	case kind == uniformKind && isBucketSizeConstant:
		return generator.GenerateUniformConstantBucketSize
	case kind == uniformKind && !isBucketSizeConstant:
		return generator.GenerateUniform
	case kind == exponentialKind && isBucketSizeConstant:
		return generator.GenerateExponentialConstantBucketSize
	case kind == exponentialKind && !isBucketSizeConstant:
		return generator.GenerateExponential
	case kind == normalKind && isBucketSizeConstant:
		return generator.GenerateNormalConstantBucketSize
	case kind == normalKind && !isBucketSizeConstant:
		return generator.GenerateNormal
	case kind == v2xKind && isBucketSizeConstant:
		return generator.GenerateV2XEnvironmentalConstantBucketSize
	case kind == v2xKind && !isBucketSizeConstant:
		return generator.GenerateV2XEnvironmental
	default:
		return nil
	}
}

func setUpGenerateFlags(command *cobra.Command) {
	command.Flags().UintP(itemCountValue, "", 30, "count of items")
	command.Flags().UintP(itemSizeValue, "", 20, "maximum size of single item")
	command.Flags().UintP(bucketCountValue, "", 10, "count of buckets")
	command.Flags().UintP(bucketSizeValue, "", 50, "size of bucket (if "+
		constantBucketSizeValue+" is not specified, the value will be randomly generated from range [1 - "+
		bucketSizeValue+"])")
	command.Flags().BoolP(constantBucketSizeValue, "", false,
		"disables random bucket sizes and enables constant size for all buckets")
	command.Flags().UintP(countValue, "", 1, "specify how many files should be generated")
	command.Flags().StringP(kindValue, "", "uniform", "specify generator kind (uniform, exponential, normal, v2x)")
}
