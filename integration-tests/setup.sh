#!/bin/sh

echo "Getting the firefox driver"
wget -nc https://github.com/mozilla/geckodriver/releases/download/v0.14.0/geckodriver-v0.14.0-linux64.tar.gz
tar -xvzf geckodriver-v0.14.0-linux64.tar.gz

echo "Getting selenium"
wget -nc https://selenium-release.storage.googleapis.com/3.0/selenium-server-standalone-3.0.1.jar

