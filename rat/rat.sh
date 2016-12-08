#!/usr/bin/env bash
wget --post-data="user=me&host=localhost&message=this is a test" http://$RATSERVER:$RATPORT/

