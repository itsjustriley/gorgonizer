[![Go Report Card](https://goreportcard.com/badge/github.com/itsjustriley/gorgonizer)](https://goreportcard.com/report/github.com/itsjustriley/gorgonizer)

# Gorgonizer
Gorgonizer is an organizer utility written in Go, designed to organize messy directories into tidy folders by type (video, document, audio, etc.). You can also have it organize by specific file type.


### Installation & How To Run 
Download or Clone this repo, then:

```
cd Gorgonizer
go mod tidy
go build -o gorgonizer
./gorgonizer --dir /path/to/directory
``` 

Note: You must include a directory when running. 


### Flags
##### `--include-subfolders`

Include this to organize subfolders. 

#### `--verbose`

Include this to print details about organization as the program runs.

#### `--defer-output`

Include this to output details about the organization process after the program is done running.

#### `--exact`

By default, Gorgonizer groups by class (eg. Images, Videos, etc. See Supported Types below). If you want it to group by specific file types, use this flag.

##### `--log`

Generates a log as the utility runs documenting action taken.

#### `--no-color` 

Disable colour in terminal output.


### Examples
#### Basic Usage
`./gorgonizer --dir ~/Downloads`

#### With subfolders and logging
`./gorgonizer --dir ~/Downloads --include-subfolders --log`

#### Exact file type organization
`./gorgonizer --dir ~/Downloads --exact`