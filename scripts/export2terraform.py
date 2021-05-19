#!/usr/bin/env python

# I know its ugly but it does the job for now.

import base64
from jinja2 import Environment, FileSystemLoader, Template
import json
import re
import argparse
import sys

rdme = """# {{ data.name }}
{{ data.description | default('', true)  }}
"""

tvar = """variable "team_id" {
    type = number
    default = "1"
}
"""

def add_receiver(data):
  for i, agent in enumerate(data["agents"]):
      data['agents'][i]['receiver_ids'] = []
  return data

def format_action(data):
  if data['diagram_layout']:
    diagram_layout = json.loads(data['diagram_layout'])
  for i, agent in enumerate(data["agents"]):
      for link in data['links']:
          if link['receiver'] == i:
              link_name = "tines_agent.{}.id".format(agent['terraform_name'])
              data['agents'][link['source']]['receiver_ids'].append(link_name)
      positions = diagram_layout.get(agent['guid'])
      position = {
        'x': positions[0],
        'y': positions[1]
      }
      data['agents'][i]['position'] = position
  return data

def format_names(data):
    for i, agent in enumerate(data["agents"]):
        data["agents"][i]["terraform_name"] = agent["name"].replace(" ", "_").lower() + "_" + str(i) 
    return data

def format_story(data):
    data["terraform_name"] = data["name"].replace(" ", "_").lower()
    return data

def get_global_resources(data):
  finds = re.findall(r'\.RESOURCE\.(.*?)[\.| |}]', str(data))
  return list(set(finds))
  
def get_credentials(data):
  finds = re.findall(r'\.CREDENTIAL\.(.*?)[\.| |}]', str(data))
  return list(set(finds))

def run(event, context):
  data = {}
  if event.get('encoded'):
    working_data = json.loads(base64.b64decode(event['export']).decode('utf-8'))
  else:
    working_data = event['export']
  format_names_data = format_names(working_data)
  receiver_data = add_receiver(format_names_data)
  format_action_data = format_action(receiver_data)
  export_data = format_story(format_action_data)
  export_data['global_resources'] = get_global_resources(export_data)
  export_data['credentials'] = get_credentials(export_data)
  env = Environment(loader=FileSystemLoader('.'))
  template = env.get_template('tines.j2')

  o = template.render(data=export_data)

  ## Write Readme
  if event.get('blog'):
    rdme_template = Template(rdme_blog)
  else:
    rdme_template = Template(rdme)
  data['readme'] = base64.b64encode(bytes(rdme_template.render(data=export_data), "utf-8")).decode('utf-8')

  ## Write main.tf
  data['main'] = base64.b64encode(bytes(o, "utf-8")).decode('utf-8')

  ## Write variables.tf
  data['vars'] = base64.b64encode(bytes(tvar, "utf-8")).decode('utf-8')

  return data

if __name__ == "__main__":
  if sys.version_info[0] != 3:
    print("This script requires Python 3")
    sys.exit(1)
  parser = argparse.ArgumentParser()
  parser.add_argument('-f', '--file', dest="file", help="Export file")
  parser.add_argument('-o', '--output', dest="output", help="Output directory")
  args = parser.parse_args()
  export = args.file
  output = args.output
  data = {}
  with open(export, 'r') as e:
    export_data = json.load(e)
  data['export'] = export_data
  files = run(data, "")
  with open(output+'README.md', 'w') as readme:
    readme.write(base64.b64decode(files['readme']).decode('utf-8'))
  with open(output+'main.tf', 'w') as tfmain:
    tfmain.write(base64.b64decode(files['main']).decode('utf-8'))
  with open(output+'variables.tf', 'w') as tfvars:
    tfvars.write(base64.b64decode(files['vars']).decode('utf-8'))
