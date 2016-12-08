#!/usr/bin/env python

import os
import socket

import requests

server = os.environ.get("RATSERVER", "localhost")
port = os.environ.get("RATPORT", "8000")

username = os.environ.get("USER", "")
hostname = socket.gethostname()

data = {
    "user": username,
    "host": hostname,
} 

def rat(message):
    data["message"] = message
    requests.post("http://{0}:{1}".format(server, port), data=data)
