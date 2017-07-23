# Notes for 23 July 2017

## Overall

Today was a successful day.  I added the longevity feature which I will now need to graph.  There are a few weird things about this - first, we are showing plenty of users who have made a single comment and then never come back.  That's probably normal, so I may add a feature to ensure we throw those out.

I may also graph the average number of users in each bucket.  I should probably do that in JS since it's a smaller data aggregation that doesn't require custom Golang work.  It takes more time to write Go than JS.

Next, testing continues.  I need to write more tests in general - my promise to myself is to keep touching files and adding tests.  That way I can run them and at least get some positive feedback.  

One thing that's tough about testing is mocking the filesystem.  I know you're not supposed to mock much in Go but this seems like a good compromise to make - it's hard to expect that a function works when it depends on the variability of a filesystem and you mock that out.

Finally, I have a lot more work to do:

- [ ] Ensure the configuration is fully integrated
- [ ] Update documentation, maybe even start writing it
- [ ] More tests
- [ ] More features and better control over them
