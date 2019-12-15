package wordcount

import "strings"

// Clone retrieves the source code from a given URL. It access the repository, clones it,
// filters non-go files and returns a channel of code.File elements.
func Clone(url string, cloner Cloner) (*Repository, <-chan File, error) {
	repo, err := cloner.Clone(url)
	if err != nil {
		return nil, nil, err
	}

	files, err := cloner.Filenames()
	if err != nil {
		return nil, nil, err
	}

	namesc := make(chan string)
	go func() {
		for _, f := range files {
			if !strings.HasSuffix(f, ".go") {
				continue
			}

			namesc <- f
		}

		close(namesc)
	}()

	filesc := make(chan File)
	go func() {
		for n := range namesc {
			rawFile, err := cloner.File(n)

			file := File{
				Name:  n,
				Raw:   rawFile,
				Error: err,
			}
			filesc <- file
		}

		close(filesc)
	}()

	return &repo, filesc, nil
}
