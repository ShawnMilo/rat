# rat

This is yet another tool written by the author, for the author. You're free to use it if you find it to be interesting.

The general idea is to log from anywhere -- multiple different applications written in different languages, running in different places, regardless of whether it's in a VM, container, or local machine.

The idea came from logging via a Slack bot. I was able to debug an application which was doing things in multiple places: A bash script on my local machine, bash and Python scripts running on a remote VM, and other code running inside a Docker container on the remote VM. This allows logging from all places to appear in a central locatoin for easier debugging.
