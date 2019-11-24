# freqtable
**freqtable** is a Frequency Table builder, which handles the extraction and counting of words (dictionary words) from Go source code.

## Input
The input for the process is a list of GitHub's public Golang source code repositories.
Since we don't use authentication to communicate to GitHub's API, we only use public available repositories.

## Output
The output after the frequency table builder execution is a plain-text file (comma-separated value, by default), which contains the following information per row:

- repository
- word
- count

The resulting file may contain a high number of rows, but it will allow us to determine and analyze each project's impact.
