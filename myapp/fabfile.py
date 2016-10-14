from fabric.api import *
import boto3
import os
import json

env.user = 'ec2-user'
env.key_filename = '~/.aws/aws_key.pem'

@task
def aws():
    session = boto3.session.Session(profile_name='default')
    client = session.client('ec2')
    instances = client.describe_instances(Filters=[
                                        {'Name':'tag-key', 'Values':['Role']},
                                        {'Name':'tag-value', 'Values':['Web']}
                                        ])
    for i in instances["Reservations"]:
          env.hosts.append(i["Instances"][0]["PublicDnsName"])

@runs_once
def build():
    local('GOOS=linux GOARCH=amd64 go build')

@task
def deploy():
    execute(build)
    sudo('sudo /usr/local/bin/supervisorctl stop myapp')
    put('myapp', '/home/ec2-user')
    sudo('sudo /usr/local/bin/supervisorctl start myapp')
