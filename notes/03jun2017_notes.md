# Multiple Executables

I realized that the big analysis file I was working with was too much and have begun to split things out into separate modules.  First, I want a module that just converts to and from proto - I should have done this earlier. Now it's done.

In the refactor, maybe even in the writing, I was calling `flush()` about a zillion times too many.  This might have thrown me off (did it write a bunch of gzip junk to the file?) but I will check later.

For now though I think I am ready to convert to proto and unmarshal successfully.
