import json

redditors = set()

badlines = 0
lines = 0

with open('./the_donald_comments.json', 'r') as infile:
    for line in infile:
        lines += 1
        try:
            data = json.loads(line)
        except ValueError:
            badlines += 1
            continue
        redditors.add(data['author'])

print('Bad lines: {} and Total Lines: {}'.format(badlines, lines))
redditors = list(redditors)
with open('./the_donald_users.json', 'w') as outfile:
    json.dump(redditors, outfile, sort_keys=True, indent=3, separators=(',', ':'))
