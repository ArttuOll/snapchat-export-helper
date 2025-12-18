package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

	counter := 0
	for node := range tree.Descendants() {
		if node.Type == html.ElementNode && node.DataAtom == atom.A {
			for _, attribute := range node.Attr {
				if attribute.Key == "onclick" {
					url, err := parseUrl(attribute.Val)
					if err != nil {
						return err
					}

					err = downloadFile(url, fmt.Sprintf("snapchat_memory_%v", counter))
					if err != nil {
						return err
					}

					counter++

					// Only download the first file during development
					if counter > 1 {
						return nil
					}
				}
			}
		}
	}

	return nil
}

func downloadFile(url *url.URL, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create a new file: %w", err)
	}
	defer file.Close()

	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to construct HTTP GET request: %w", err)
	}
	request.Header.Add("X-Snap-Route-Tag", "mem-dmd")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		os.Remove(filename)
		return fmt.Errorf("failed to download file %v: %v", url.String(), err)
	}
	defer response.Body.Close()

	if response.StatusCode >= 299 {
		responseMessage, _ := io.ReadAll(response.Body)
		return fmt.Errorf("snapchat's server returned with non-success status code: %v %v", response.StatusCode, string(responseMessage))
	}

	log.Printf("server responded with status %v\n", response.StatusCode)

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response to output file: %v", err)
	}

	return nil
}

func parseUrl(input string) (*url.URL, error) {
	urlString := strings.Split(input, "'")[1]
	url, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL from %s", urlString)
	}

	return url, nil
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
