# -*- coding: utf-8 -*-

"""
Given a map of users -> subreddits -> comment counts

We want to find which of the reddits are most common.

* Note: We should maybe just keep misc as a category?
"""

from operator import itemgetter
from collections import defaultdict

import json

num_features = 100
infile = './test_profile_data/donald_users_with_subreddits.json'

with open(infile, 'r') as input_file:
    data = json.load(input_file)

reddit_counts = defaultdict(int)

for key, value in data.items():
    for sub in value:
        reddit_counts[sub.lower()] += 1

# We do not need this one.
del reddit_counts['the_donald']

comment_counts = []

for key, val in reddit_counts.items():
    comment_counts.append([key, val])

comment_counts.sort(key=itemgetter(1))

key_features = comment_counts[-1 * num_features:]
print(json.dumps(key_features, indent=3, separators=(',', ': ')))

# Calculate average number of subs per user.
subs_per_user = len(data) / len(reddit_counts)
print('Number of subs per user: {}'.format(subs_per_user))

tgt_features = [_[0] for _ in key_features]


def generate_feature_array(target_features):
    """Take a list of subs and generate a dict of x->0 for every item."""
    ret_val = dict()
    for item in target_features:
        ret_val[item] = 0
    return ret_val


def gen_feature_vector(in_dict):
    """Take a dict of item->val, sort it and return a set of [x, x, x, x...]"""
    values = sorted(in_dict.items(), key=itemgetter(0))
    vector = [_[1] for _ in values]
    return vector


def normalize_vector(in_features):
    """Normalize things."""
    pass


def generate_features_from_users(user_data, target_features):
    features = []
    for user, values in user_data.items():
        features = generate_feature_array(target_features)
        for sub, comments in values.items():
            features[sub] = comments
        features.append(gen_feature_vector(features))


def test_things():
    """
    take the features
    generate a classifier
    classify everything
    """
    pass

