package cmd

import (
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type fileIoer interface {
	GetOutputFileName(afero.Fs, string) (string, error)
}

type contactCleaner interface {
	ContactClean(afero.Fs, string, string)
}

type CommandError struct {
	Msg string
}

func (err CommandError) Error() string {
	return err.Msg
}

func RootCmd(fileIo fileIoer, contactClean contactCleaner) (rootCmd *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "cleancontacts <filepath>.vcf",
		Short: "Clean Contacts removes fields from your phone contacts to prevent apps for having access more than they need.",
		Long: `Do not want to share all your contact info to mobile apps? This software is for you!

Export your contacts in VCard format and run the program. BANG! You have a new VCard file with cleaned contacts with just their names and telephones.`,
		Version: "0.2.2",
		Example: `cleancontacts contacts.vcf
cleancontacts /path/contacts.vcf`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return CommandError{Msg: "Contact file argument not provided."}
			}
			filePath := args[0]

			dir, fileName := filepath.Split(filePath)
			fsys := afero.NewOsFs()
			fileNameOut, err := fileIo.GetOutputFileName(fsys, fileName)
			if err != nil {
				return err
			}
			filePathOut := filepath.Join(dir, fileNameOut)

			contactClean.ContactClean(fsys, fileName, filePathOut)
			return nil
		},
	}
	return cmd
}

func Execute(cmd *cobra.Command) error {
	return cmd.Execute()
}
