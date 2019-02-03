# moni

[![CircleCI](https://circleci.com/gh/adrian-gheorghe/moni.svg?style=svg)](https://circleci.com/gh/adrian-gheorghe/moni)

moni (short for monitoring) is a utility that scans your file system periodically and alerts you when your file signatures have changed. Can be configured to execute different commands on failure or success.

## Download
Download latest from the releases page: https://github.com/adrian-gheorghe/moni/releases

## Configuration

```yaml
general:
  # Should moni keep running and execute periodically  
  periodic: true
  # If periodic is true, what interval should moni run at? Interval value is in seconds
  interval: 50
  # Tree is stored as a json to the following path
  tree_store: ./output.json
  # Path to parse
  path: /var/www/html
  # Command that is run if the tree is identical to the previous one
  command_success: "echo SUCCESS"
  # Command that is run if the tree is not identical to the previous one
  command_failure: "echo FAILURE"
log:
  # Log path for moni  
  log_path: ./log.log
  # Memory log options are only for development use. Please keep memory_log value to false
  memory_log_path: ./memory.log
  memory_log: false
algorithm:
  # Algorithm options are:
  # - FlatTreeWalk (manual recursive treewalk)  
  # - GoDirTreeWalk - walk algorithm developed by karrick - https://github.com/karrick/godirwalk
  name: FlatTreeWalk
  processor: ObjectProcessor
  # List of directory / file names moni should ignore
  ignore:
    - ".git"
    - ".idea"
    - ".vscode"
    - ".DS_Store"
    - "node_modules"
    - "uploads"
```

## Usage

```bash
./moni --help
Usage of ./moni:
  -config string
    	path for the configuration file (default "./config.yml")
  -version
    	Prints current version
```
Run moni with the config flag pointing to the path to your configuration yml file. The default config path is config.yml in the current directory

```bash
./moni --config="./config.yml"
```

## Docker Configuration
You can use the official image from Docker Hub to raise a container that monitors files in a volume. https://cloud.docker.com/repository/docker/adighe/moni . You can change the configuration by changing the environment variable CONFIG_PATH to your config.ym file. The configuration embedded in the container is the default one from sample.docker.config.yml

```yaml
version: '3.3'
services:
  moni:
    image: adighe/moni:latest
    volumes:
      - ./:/var/www/html
    environment:
      CONFIG_PATH: /app/config.yml   
```