# aoc2md
CLI to help initialise the solution for Advent of Code 🎄. Solving it though, it's on you 😉

## Install

I maintain a homebrew-tap repo with my CLI's. You can just install it with Brew or check the releases pages and install it that way.

```sh
brew install kevinrobayna/homebrew-tap/aoc2md
```

## Session Value
Before you run the program you need to grab the 🍪 cookie session value from [adventofcode.com](adventofcode.com) after you log in.

You should do this every time the Advent of Code season starts since the session will expire.

## Usage

To generate today's problem just run the following:

```sh
aoc2md
```

But you can specify the day and year like so:

```sh
aoc2md --day 1 --year 2015
```

Remember to make the session available through the env variable `AOC_SESSION` or directly in the command like so:

```sh
aoc2md --session_id your_secret_value
```
