# Collate

## Installation
`go get -u github.com/whytheplatypus/collate`

## Usage
`collate -template=<tempalte_file> -token=<github_token> <org> <repo> <issue number>`

If `-token` is not given then `collate` will look for one in `~/.ghtoken`.

You can see defaults and flag usage with `collate -h`
