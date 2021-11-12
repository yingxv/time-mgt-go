#!/bin/bash
set -e

tag=ngekaworu/time-mgt-go

docker build --file ./Dockerfile --tag ${tag} ..;
docker push ${tag};
