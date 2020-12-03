#!/usr/bin/env python

import json

with open('export.json', 'r') as f:
    data = json.load(f)

for i, agent in enumerate(data["agents"]):
    a = """resource "tines_agent" "{0}" {{
    name = "{0}"
    agent_type = "{1}"
    story_id = var.story_id
    keep_events_for = {2}
    source_ids = []
    receiver_ids = []
    agent_options = jsonencode({3})
}}""".format(agent['name'].replace(' ', '_'), agent['type'], agent['keep_events_for'], json.dumps(agent['options']))
    print(a)
