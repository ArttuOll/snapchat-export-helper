package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download all the export files",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: download,
	Args: cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}

func download(cmd *cobra.Command, args []string) error {
	path := args[0]
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open %v. is that the correct path to the memories_history.html file?", path)
	}

	s := string(b)

	_, err = html.Parse(strings.NewReader(s))
	if err != nil {
		return fmt.Errorf("failed to parse momeries_history.html file")
	}

	return nil
}
