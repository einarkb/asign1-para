# IMT2681 Assignment 2 - Paragliding
## Created by Einar Budsted - 131348

I were not figure out how to deploy the clock_trigger on openstack. instead I tested it locally up against the api on heroku and it worked great.
I Also used Discord webhooks instead of Slack.

The genreal architecture of the code ended up being a little more coupled than intended, and due to problems deploying on heroku I ended
up moving most of the code into a single package which made sense anyway whith how coupled the code has become.

