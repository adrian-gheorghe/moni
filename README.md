# moni

[![CircleCI](https://circleci.com/gh/adrian-gheorghe/moni.svg?style=svg)](https://circleci.com/gh/adrian-gheorghe/moni)

moni (short for monitoring) is a utility written in go that scans your file system periodically and alerts you when your file signatures have changed. Can be configured to execute different commands on failure or success.

## Download
Download latest from the releases page: https://github.com/adrian-gheorghe/moni/releases

## Docker Configuration
You can use the image from Docker Hub to raise a container that monitors files in a volume. https://hub.docker.com/r/adighe/moni/ . You can change the configuration by changing the environment variable CONFIG_PATH to your config.yml file and mounting your file to the container. The configuration embedded in the container is the default one from sample.docker.config.yml. Logs are sent to stdout for this configuration

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
  # Log path for moni. Accepts file path or "stdout"
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

## Sample output from Docker container

```bash
moni_1  | flat_treewalk.go:25: File count:  16
moni_1  | processor.go:53: Tree has changed
moni_1  | processor.go:54: {main.TreeFile}.Children:
moni_1  | 	-: []main.TreeFile{{Path: "/var/www/html", Type: "directory", Mode: "drwxr-xr-x", Size: 256, Modtime: "2019-02-03 15:28:56 +0000 UTC"}, {Path: "/var/www/html/a", Type: "directory", Mode: "drwxr-xr-x", Size: 160, Modtime: "2019-01-27 16:45:34 +0000 UTC"}, {Path: "/var/www/html/a/ac", Type: "directory", Mode: "drwxr-xr-x", Size: 96, Modtime: "2019-01-27 16:45:33 +0000 UTC"}, {Path: "/var/www/html/a/ac/acd.txt", Type: "file", Mode: "-rw-r--r--", Modtime: "2019-01-27 16:45:33 +0000 UTC", Sum: "d41d8cd98f00b204e9800998ecf8427e"}, {Path: "/var/www/html/a/ab.txt", Type: "file", Mode: "-rw-r--r--", Size: 22, Modtime: "2019-01-27 16:45:34 +0000 UTC", Sum: "8715dae36a0b112120136e6e52258063"}, {Path: "/var/www/html/a/az.txt", Type: "file", Mode: "-rw-r--r--", Size: 4, Modtime: "2019-01-27 16:45:33 +0000 UTC", Sum: "098f6bcd4621d373cade4e832627b4f6"}, {Path: "/var/www/html/c.txt", Type: "file", Mode: "-rw-r--r--", Modtime: "2019-01-27 16:45:33 +0000 UTC", Sum: "d41d8cd98f00b204e9800998ecf8427e"}, {Path: "/var/www/html/VERSION", Type: "file", Mode: "-rw-r--r--", Size: 5, Modtime: "2019-02-03 12:52:07 +0000 UTC", Sum: "872ccd9c6dce18ce6ea4d5106540f089"}, {Path: "/var/www/html/d.txt", Type: "file", Mode: "-rw-r--r--", Size: 12, Modtime: "2019-02-03 15:29:15 +0000 UTC", Sum: "2d486a3582cc9354f668d42cab28525f"}, {Path: "/var/www/html/config.yml", Type: "file", Mode: "-rw-r--r--", Size: 437, Modtime: "2019-02-03 13:11:20 +0000 UTC", Sum: "c180483e726af54ca90010f7e829d730"}, {Path: "/var/www/html/b", Type: "directory", Mode: "drwxr-xr-x", Size: 160, Modtime: "2019-01-27 16:45:33 +0000 UTC"}, {Path: "/var/www/html/b/bc", Type: "directory", Mode: "drwxr-xr-x", Size: 96, Modtime: "2019-01-27 16:45:33 +0000 UTC"}, {Path: "/var/www/html/b/bc/bcd.txt", Type: "file", Mode: "-rw-r--r--", Modtime: "2019-01-27 16:45:33 +0000 UTC", Sum: "d41d8cd98f00b204e9800998ecf8427e"}, {Path: "/var/www/html/b/ba.txt", Type: "file", Mode: "-rw-r--r--", Modtime: "2019-01-27 16:45:33 +0000 UTC", Sum: "d41d8cd98f00b204e9800998ecf8427e"}, {Path: "/var/www/html/b/ba", Type: "directory", Mode: "drwxr-xr-x", Size: 96, Modtime: "2019-01-27 16:45:33 +0000 UTC"}, {Path: "/var/www/html/b/ba/bdf.txt", Type: "file", Mode: "-rw-r--r--", Modtime: "2019-01-27 16:45:33 +0000 UTC", Sum: "d41d8cd98f00b204e9800998ecf8427e"}}
moni_1  | 	+: []main.TreeFile(nil)
moni_1  | 
moni_1  | processor.go:109: FAILURE
```
