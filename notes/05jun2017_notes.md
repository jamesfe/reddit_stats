# Notes for protobuf delimiting

 I learned a few things about my “easy win” in converting things to protobuf today:
 - I have spent more time on this than I originally thought I would
 - Experiential Problems:
    - I encountered a problem by which I used to delimit things with ‘\n’ but changed it to a byte that was utilized.
    - I originally (foolishly) thought “maybe I will just change the delimiter to be a more specific pattern” but that was dumb.
    - After finding that to be troublesome (and risky) I searched around the web and found a better solution.
    - Someone else thought this solution up and presented it to the world, which was really nice, without this sharing I would likely never have thought up the solution.
 - I am really beginning to learn that a lot of the time you spend in a project like this is not writing analysis code but formatting your data so it is easy to analyze.

# Results

I'm also coming to the slghtly disheartening idea that maybe protobuf is not faster than JSON in this case.  Maybe my hard drive can support the reads and it's actually the processing time that's doing it.

I'll do a few things here:
    - Make sure I run the tests when I'm not doing anything else on the computer
    - Run larger and larger sets to see if there is some gain to be had over large amounts of data vs small
    - Check my code with a profiler to see if I am being wasteful in some way

However, there is always a chance that it is taking more time to process the protobufs (i.e. decode them from bytes on disk) than it is worth reading them.  If this is true then I have probably wasted a lot of time in writing scripts to convert this data into protobuf form, but regardless, I think it was a great learning experience.
