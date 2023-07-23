# godeps

- [x] backup old go.mod
- [x] go get ./...
- [x] go mod tidy
- [x] git diff go.mod
- [x] parse diff
- [x] iterate over changes
  - [x] check if there is an existing PR
  - [x] if not
    - [x] create new branch
    - [x] patch the single line
    - [x] commit
    - [x] push
    - [x] create new PR
