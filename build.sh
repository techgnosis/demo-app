#! /usr/bin/env bash

set -euo pipefail

pack build harbor.lab.home/library/demo-app:latest

docker push harbor.lab.home/library/demo-app:latest