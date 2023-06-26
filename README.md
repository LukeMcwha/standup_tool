# Workplace Standup Tool

A basic standup tool that records information based on the configuration in 'config.json'.

The names within the config file will be shuffled to determine the order of the day. 

All yesterdays responses will be recorded. If you leave a field blank and press enter, it will use yesterdays results.

## What you need.

You will need to rename the `example_config.json` to `config.json` as the program has that file name hard coded.

## How to build

`go build .` builds a binary of the program. In my case it calls it `standup`

## How to run

Run a terminal at the binary location. Run `./standup` to run the executable.
