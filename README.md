# godeps

- [x] backup old go.mod
- [x] go get ./...
- [x] go mod tidy
- [ ] git diff go.mod
- [x] parse diff
- [ ] iterate over changes
  - [ ] check if there is an existing PR
  - [ ] if not
    - [ ] create new branch
    - [ ] patch the single line
    - [ ] commit
    - [ ] push
    - [ ] create new PR
