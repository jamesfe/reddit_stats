# -*- coding: utf-8 -*-

import json
from collections import defaultdict
import os
# import gzip

weird_subs = defaultdict(int)
with open('./the_donald_users.json', 'r') as userfile:
    target_authors = json.load(userfile)
authors = set()


def decode_line(line):
    try:
        # data = json.loads(line.decode('utf-8'))
        data = json.loads(line)
        return data
    except Exception as e:
        print('Bad Line: {} with error {}'.format(line, e))


def author_func(line):
    data = decode_line(line)
    if data['author'] in target_authors:
        weird_subs[data['subreddit']] += 1


def just_get_author(line):
    data = decode_line(line)
    if data:
        authors.add(data['author'])


def for_every_file_exec(datadir, func):
    i = 0
    for item in os.listdir(datadir):
        i += 1
        if i > 2:
            print('Dropping out of file loop')
            break
        if item.lower()[-2:] == 'gz':
            print(item)
            with open(os.path.join(datadir, item), 'r') as k:
                for line in k:
                    func(line)


def main():
    # input_dir = '/Users/jferrara/PersCode/reddit_donald/data'

    filtered_dir = '/Users/jferrara/PersCode/reddit_stats/filters/old'
    print('Looping')
    for_every_file_exec(filtered_dir, author_func)
    print('Done with file IO')

    print('Logging')
    with open('./the_donald_users_2.json', 'w') as outfile:
        json.dump(list(authors), outfile, sort_keys=True, indent=3, separators=(',', ':'))
    print('Done logging')
    # with open('./the_donald_sub_count.json', 'w') as outfile:
    #    json.dump(weird_subs, outfile, sort_keys=True, indent=3, separators=(',', ':'))


main()
