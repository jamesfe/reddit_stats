# Random Sample

We need to randomly sample all the subreddits to figure out if there are a lot of deleted comments or not.  For this I think I will sample about .25% of all the comments which should be roughly 325,000,000 comments.  I'll do this by adding a new function similar to "isDonaldLite" and just dump that into my functional check, then I'll refactor later.

# Lessons Learned

You may think it's a good idea to use some random numbers when you sample.  I agree.  Here's a bad idea:

```
func IsRandomSample(percentTrue float32) bool {
	/* given the percent value we are given, return true if that percent hits randomly. */
	rand.Seed(time.Now().UTC().UnixNano())
	val := rand.Float32()
	if val*100 < percentTrue {
		return true
	} else {
		return false
	}
}
```

Seeding the random number generator **every time**!?  Are you crazy?  Yes, I was.  What a horrid mistake.  Don't do this.  I profiled the program before and after doing this and it turns out the CPU spend 8.85s/56.99s just seeding the random number generator. 

You can see the profiles for this here:

[profiling just checking if it's an /r/the_donald post](./profiles/checking_reddit_is_the_donald_profile.svg)

[profiling while seeding every time](.profiles/checking_random_sample_profile.svg)

[profiling with refactored seeding](./profiles/checking_randoms_seeded_once_profile.svg)
