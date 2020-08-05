# watch output of `defaults read`

# TODO
- [x] Diff not working when reading from defaults read
- [x] When outputting a `defaults` command, the Marshaled string does not escape newlines and actually prints them out
    - Does work when you copy the entire command and run it.
- [x] Create a loop to poll `defaults read`

- [ ] Code clean up
- [ ] Tests
    - [ ] Test that this works with data
- [ ] clean up go-plist
  - [ ] Move into own repo
  - [ ] Gut excess functionality unnecessary for `defaults`
