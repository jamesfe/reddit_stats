# Filtering

I finally decided to make a program that filters on a subreddit and outputs all its comments to another directory.  This was pretty easy - I borrowed a lot of code from the "analysis" portion.

## Polymorphism

One thing I'd like to do is make this more polymorphic (is this the right word?) in that the file-reading function could take a function that performs analysis and returns a varying type of result but this seems a little much and reflection is pretty slow in Golang IIRC.

## Gzipping Failures

Another thing I would like to figure out is why my gzipped files won't ungzip.  You may notice in this commit (6dc91a0d0de38bdb385c38d0fb37b8b2d48533c4) I add a bunch of logging and stuff for the GetFileWriter function but none of it seems to trigger and I still can't unzip gzips!  The error message is something like "Unexpected EOF".  Bummer.
