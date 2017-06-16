#!/bin/bash
set -ev
if [ "${TRAVIS_PULL_REQUEST}" = "false" ]; then
  curl 'https://goreportcard.com/checks' -H 'Accept: application/json, text/javascript, */*; q=0.01' -H 'Accept-Language: en-US,en;q=0.5' --compressed -H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' -H 'X-Requested-With: XMLHttpRequest' -H 'Referer: https://goreportcard.com/report/github.com/littledot/mockhiato' -H 'Connection: keep-alive' --data 'repo=github.com%2Flittledot%2Fmockhiato'
fi
