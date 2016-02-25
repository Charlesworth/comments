# comments

A comments service for static sites, WIP branch

## what should this thing look like:
#### inputs
- POST from a Form (use request.ParseForm() then act on the r.Form to access)

#### outputs
- GET /:url (JSON)
- GET /:url?callback=whatever (JSONP)
- FUTURE: GET /:url?withScript=true (js with comment HTML and JSON data embedded) !!check if having POST form data will make a diff to using query string

#### persistant storage
- byteStore (boltDB backend) DONE

#### user config
- set port DONE
- FUTURE: add .js file that ?withScript=true will embed it
- FUTURE: backup backend support like aws s3 and google cloud storage
