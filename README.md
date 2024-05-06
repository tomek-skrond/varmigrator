# varmigrator
Application for editing and viewing Github Actions secrets/variables written in Go.

# Building the app
To build the app, clone the repository, then in the working directory issue a command:
```
$ go build
```
This command will output an executable `varmigrator`.

# Demo
To start using `varmigrator` export your github `personal access token` an environment variable `GITHUB_TOKEN`.<br>
Token should have read/write permissions to repositories you want to edit/view. <br>
Additionaly provide `repository` and `username` you want to view.
```
$ ./varmigrator -repo microservice-infra -username tomek-skrond
All available project variables:
1. GINGER (current value: ginger)

All available project secrets:
1. FINGER

Which variable do you want to edit?
Enter in a format vN or sN (variable/secret number N)
Your input: s1
```
In this case, we enter `s1`. This means that we edit first secret in the list - `FINGER`. If you want to edit a regular variable `GINGER` - enter `v1`.

To only view all variables/secrets without editing, we can use `-print` option, also we can modify output with additional flags. Example below:
```
$  ./varmigrator -repo microservice-infra \
                 -username tomek-skrond \
                 -print \
                 -mode json \
                 -pretty \
                 -concise
[
  {
    "Id": 0,
    "name": "GINGER",
    "value": "ginger",
    "created_at": "2024-05-05T23:53:54Z",
    "updated_at": "2024-05-06T17:32:30Z"
  }
]

[
  {
    "Id": 0,
    "name": "FINGER",
    "created_at": "2024-05-05T23:53:36Z",
    "updated_at": "2024-05-06T17:32:40Z"
  }
]
```

All functionalities:
```
$ ./varmigrator -help
Varmigrator - tool for editing and viewing Github Repository Secrets/Variables
Options:
  -concise
        Only print data
  -help
        Print usage information
  -mode string
        Print mode (normal/json) (default "normal")
  -pretty
        Pretty-print json
  -print
        Print-only mode
  -repo string
        GitHub Repo name
  -secrets-only
        Print only repository variables
  -username string
        GitHub Repo owner
  -vars-only
        Print only repository secrets
  -verbose
        Verbose mode
```
