# Cleanup and Direction

The last few days have involved a lot of cleanup: I learned significant lessons about my journey and I've implemented a few of them.

First, I removed protobuf support.  I created a tag in the GitHub repo to mark this moment in time.  Pour one out for protobuf, it just wasn't right for this project.  Maybe I will have the heart to try another serialization format later, but today I am not interested in expending the effort again, although I learned a lot from it.
Second, I refactored a bunch of stuff into this new directory structure which is much neater.  That's nice.

Next I need to delete some unused code since it's just cruft hanging around.  I can go back in git and get it later if I want it back.

Finally, my next anaysis will be: deleted vs non-deleted comments. Deleted comments have author='[deleted]'.  This will be easy and then I will go back over to d3 and reuse some code to graph it out - maybe a nice measure of censorship per reddit?

I also need to think about two other things:
1. How to run multiple analyses concurrently in the same input loop (maybe another tack into goroutines)
2. How to run multiple input loops easily (a refactoring job mainly)

