#!/bin/bash

curl -X GET 'https://api.twitch.tv/helix/eventsub/subscriptions' -H 'Authorization: Bearer TOKEN' -H 'Client-Id: CLIENT_ID'