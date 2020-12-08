#!/usr/bin/env python

# I know its ugly but it does the job for now.


from jinja2 import Environment, FileSystemLoader, Template
import json
import argparse
from os.path import abspath

rdme = """# {{ data.name }}
{{ data.description }}
"""

tvar = """variable "tines_token" {
    type = string
}

variable "tines_email" {
    type = string
}

variable "tines_base_url" {
    type = string
}

variable "story_id" {
    type = number
}
"""

def add_receiver(data):
  for i, agent in enumerate(data["agents"]):
      data['agents'][i]['receiver_ids'] = []
  return data

def format_receiver(data):
  for i, agent in enumerate(data["agents"]):
      for link in data['links']:
          if link['receiver'] == i:
              link_name = "tines_agent.{}.id".format(agent['name'].replace(' ', '_').lower())
              data['agents'][link['source']]['receiver_ids'].append(link_name)
  return data

def run(export, output, readme):
    
  if sys.version_info[0] != 3:
    print("This script requires Python 3")
    sys.exit(1)
  with open(export, 'r') as e:
    export_data = json.load(e)

  export_data = add_receiver(export_data)
  export_data = format_receiver(export_data)
  env = Environment(loader=FileSystemLoader('resources/'))
  template = env.get_template('tines.j2')

  o = template.render(data=export_data)

  if output:
    path = abspath(output)
    if not path.endswith('/'):
      path += "/"
    if readme:
      with open(path+'README.md', 'w') as readme:
        rdme_template = Template(rdme)
        readme.write(rdme_template.render(data=export_data))
    with open(path+'main.tf', 'w') as tfmain:
      tfmain.write(o)
    with open(path+'variables.tf', 'w') as tfvars:
      tfvars.write(tvar)

if __name__ == "__main__":
  parser = argparse.ArgumentParser()
  parser.add_argument('-f', '--file', dest="file", help="Export file")
  parser.add_argument('-o', '--output', dest="output", help="Output directory")
  parser.add_argument('-r', '--readme', dest="readme", help="Create a README file for the story")
  args = parser.parse_args()
  export = args.file
  output = args.output
  readme = args.readme
  run(export, output, readme)
