#!/usr/bin/env python

import json

with open('export.json', 'r') as e:
    data = json.load(e)

for i, agent in enumerate(data["agents"]):
    data['agents'][i]['receiver_ids'] = []

for i, agent in enumerate(data["agents"]):
    for link in data['links']:
        if link['receiver'] == i:
            link_name = "tines_agent.{}.id".format(agent['name'].replace(' ', '_').lower())
            data['agents'][link['source']]['receiver_ids'].append(link_name)

for i, agent in enumerate(data["agents"]):
    receiver_id_string = "[{}]".format(', '.join(agent['receiver_ids']))
    a = """resource "tines_agent" "{0}" {{
    name = "{0}"
    agent_type = "{1}"
    story_id = var.story_id
    keep_events_for = {2}
    source_ids = []
    receiver_ids = {4}
    agent_options = jsonencode({3})
}}""".format(agent['name'].replace(' ', '_').lower(), agent['type'], agent['keep_events_for'], json.dumps(agent['options']), receiver_id_string)
    print(a)
