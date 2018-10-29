# IMT2681 Assignment 2 - Paragliding
## Created by Einar Budsted - 131348

## Api made in go for paragliding

### Four environment variables are being used, one of them is optional
- PORT: The port the app is listening on
- DB_URI: the uri used to connect to the database
- DB_NAME: the name of the database
- N_TICKER_PAGE(optional): number of entries ticker reponds with for paging. if not set it will default to 5 

I were not able to figure out how to deploy the clock_trigger on openstack. instead I tested it locally up against the api on heroku and it worked great.
I Also used Discord webhooks instead of Slack as I am not fammiliar with slack and am a big faen olf discord

The genreal architecture of the code ended up being a little more coupled than intended, and due to problems deploying on heroku I ended
up moving most of the code into a single package which made sense anyway whith how coupled the code has become.

I chose to use db connection info directly in the code for tests in case teachers is planning to run them.

