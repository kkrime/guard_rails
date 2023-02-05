import requests
import sys
import json
from pygments import highlight
from pygments.lexers import JsonLexer
from pygments.formatters import TerminalFormatter
import argparse
from termcolor import colored


# PARSAE ARGUMENTS
def proccessEnv():
 global host
 global port

 parser = argparse.ArgumentParser(description='repsoitory e2e tests')
 parser.add_argument('--env',
                    dest='env',
                    type=str,
                    default='local',
                    choices=['local', 'staging', 'production'],
                    help='which env to run the test against')

 env = parser.parse_args().env

 if env == 'local':
  host ='http://localhost'
  port = 8080

 # coming soon
 # elif env == 'staging':
  # host =
  # port =

 # coming soon
 # elif env == 'production':
  # host =
  # port =

proccessEnv()

# HELPER FUNCTIONS
def addRepository(data):
 response = requests.post( host + ':' + str(port) + url, data=json.dumps(data), headers=headers)
 jsonResponse = response.json()
 printJson(jsonResponse)
 return jsonResponse

def getRepository(repositoryName):
 response = requests.get( host + ':' + str(port) + url + repositoryName, headers=headers)
 # print(response.content)
 jsonResponse = response.json()
 printJson(jsonResponse)
 return jsonResponse

def changeRepositoryUrl(data):
 response = requests.put( host + ':' + str(port) + url, data=json.dumps(data), headers=headers)
 jsonResponse = response.json()
 printJson(jsonResponse)
 return jsonResponse

def deleteRepository(repositoryName):
 response = requests.delete( host + ':' + str(port) + url + repositoryName, headers=headers)
 # print(response.content)
 jsonResponse = response.json()
 printJson(jsonResponse)
 return jsonResponse

def printJson(out):
 formattedJson = json.dumps(out, indent=2)
 print(highlight(formattedJson, JsonLexer(), TerminalFormatter())) 

def passes(data):
 if data["status"] != "success":
  sys.exit("Should pass")

def fails(data):
 if data["status"] != "success":
  sys.exit("Should fail")

# VARS
repositoryName = "repository"
repositoryUrl = "https://github.com/torvalds/linux"

headers = { 'Content-Type':'application/json' }

url = '/v1/repository/'
inputData = {}
inputData["name"] = repositoryName
inputData["url"] = repositoryUrl

# SETUP
print (colored('0. Delete repository {}'.format(repositoryName) , 'green'))
deleteRepository(repositoryName)

# TESTS
print (colored('1. Adding repository {}'.format(repositoryName), 'green'))
result = addRepository(inputData)
passes(result)

print (colored('2. Get repository {}'.format(repositoryName), 'green'))
result = getRepository(repositoryName)
passes(result)
# result  = json.load(result)

repositoryUrl = "https://github.com/golang/go"
print (colored('3. change repository url to {}'.format(repositoryUrl), 'green'))
inputData["url"] = repositoryUrl
result = changeRepositoryUrl(inputData)
passes(result)

print (colored('4. Get repository {}'.format(repositoryName) + ' with new url {}'.format(repositoryUrl), 'green'))
result = getRepository(repositoryName)
passes(result)
if result['data']['url'] != repositoryUrl:
 sys.exit(colored('url has not been changed', 'red'))

print (colored('5. Delete repository {}'.format(repositoryName) , 'green'))
result = deleteRepository(repositoryName)
passes(result)

print (colored('4. Check repository {}'.format(repositoryName) + ' is deleted', 'green'))
result = getRepository(repositoryName)
fails(result)

print (colored('All Tests Pass Successfully!', 'green'))
