# -*- coding: utf-8 -*-

import json
from collections import defaultdict
import os
import gzip

weird_subs = defaultdict(int)

with open('./the_donald_users.json', 'r') as userfile:
    target_authors = json.load(userfile)

datadir = '/Users/jferrara/PersCode/reddit_donald/data'

for item in os.listdir(datadir):
    if item.lower()[-2:] == 'gz':
        with gzip.open(os.path.join(datadir, item), 'r') as k:
            for line in k:
                try:
                    data = json.loads(line.decode('utf-8'))
                except Exception as e:
                    print('Bad Line: {} with error {}'.format(line, e))
                    continue
                if data['author'] in target_authors:
                    weird_subs[data['subreddit']] += 1

with open('./the_donald_sub_count.json', 'w') as outfile:
    json.dump(weird_subs, outfile, sort_keys=True, indent=3, separators=(',', ':'))
