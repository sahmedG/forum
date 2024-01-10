#!/bin/bash

echo "building image.."

sleep 1

docker build -t forumimage .

echo "starting server.."

sleep 1

docker run -p 8080:8080 forumimage