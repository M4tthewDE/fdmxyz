#!/bin/bash

curl -X POST http://localhost:8080/submit -d '{"user":"test", "ranking":["Germany", "England", "Italy"]}'