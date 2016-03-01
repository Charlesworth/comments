# comments

A comments service for static sites, currently work in progress. Would welcome any feature requests or issues. Current features:
- Easy to set up: a single binary file for each platform to return
- Supports JSONP and CORS
- Single database file for easy backup
- Dockerfile provided for easy deployment

## API

### - Add comment:
##### `POST` /[page]
Adds the comment to the [page] bucket. POST parameters:
- `msg` the comments message
- `poster` the poster of the comment

### - Delete comment:
##### `DELETE` /[page]/[time]
Delete the comment in [page] bucket with the unix time [time]

### - Get comments:
##### `GET` /[page]
Returns an array comments for [page] in JSON format
##### `GET` /[page]?callback=[callback]
Returns all comments for [page] in JSON format with the wrapper function name [callback]
##### return format:
```
[{
  "Poster":"Charlie",
  "Page":"blogPost1",
  "Msg":"Great blog post, can't wait to read more",
  "TimeUnix":"1456415718168122396"
 },
 {
  "Poster":"Sally",
  "Page":"blogPost1",
  "Msg":"Terrible blog post, I dissagree",
  "TimeUnix":"1456415718168154378"},   
}]
```

## Using Dockerfile
A Dockerfile is provided for use with this project. To build and run comments from a docker container:

    $ docker build -t [name_of_image] .
    $ docker run -p 8000:[port_on_host] [name_of_image]

## TODO
#### config
- set the CORS accepted origins header on startup.

#### tests
- get good (80%+) test coverage.

#### Releases
- output bins for Windows and Linux.
