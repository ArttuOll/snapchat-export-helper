package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
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
	s, err := readInputToString(args[0])
	if err != nil {
		return err
	}

	tree, err := parse(s)
	if err != nil {
		return err
	}

	for node := range tree.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			for _, attribute := range node.Attr {
				if attribute.Key == "onclick" {
					fmt.Println(attribute.Val)
				}
			}
		}
	}

	return nil
}

func readInputToString(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to open %v. is that the correct path to the memories_history.html file?", path)
	}

	return string(b), nil
}

func parse(input string) (*html.Node, error) {
	tree, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("failed to parse memories_history.html file")
	}

	return tree, nil
}
