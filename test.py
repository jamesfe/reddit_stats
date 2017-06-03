
i = 0
outfile = open('./the_donald_comments.json', 'w')
items = 0
with open('~/PersCode/reddit_donald/data/RC_2017-03', 'r') as blah:
    for item in blah:
        if item.lower().find("the_donald") > -1:
            outfile.write(item + "\n")
            items += 1
        i += 1

print("iterated: ", i)

outfile.close()
print("Wrote: ", items)
