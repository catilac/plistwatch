# watch output of `defaults read`

# TODO
- [ ] When outputting a `defaults` command, the Marshaled string does not escape newlines and actually prints them out
    - Does work when you copy the entire command and run it.
- [ ] Need to create a loop to poll `defaults read`
- [ ] Code clean up
- [ ] Tests
- [ ] clean up go-plist
  - [ ] Move into own repo
  - [ ] Gut excess functionality unnecessary for `defaults`
