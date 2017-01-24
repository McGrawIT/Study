# Slot Sevice (slot-service)
(description needed)

This sample app depends on:  
[MS Excel Ingest / Convert](https://github.com/tealeg/xlsx) For conversion of Slot Data Files in MS Excel format into internal Slot Data DB.  File format: AAA...(DDMMM-DDMMMYY).xlsx, where AAA = 3-digit Airport Code, DDMMM and DDMMMYY are the Begin and End Dates (Saason) for this Slot Data.  The ... is typically a description of the File (e.g., Summer 2016), but can be any characters, except "(".  
[Amazon Web Searcives](https://github.com/aws/aws-sdk-go/aws") Golong SDK supporting repositories (including /session, request, /credentials) for AWS access.
[AWS Simple Storage Service (S3)](https://github.com/aws/aws-sdk-go/service/s3) S3-specific Service Golang SDK repoositories (includes /s3manager).
All of these dependencies can be removed if the relevant functionality is not required. A basic web service can be written with no dependencies beyond the Go standard libraries. If you require any third party libraries, you will also require Godep for depedency vendoring.

## Pushing a Go app to Predix
(WIP) To push this app to Predix:  
* Get the app source:  
  `go get -d github.build.ge.com/AviationRecovery/slot-service.git`  
* Get and install Godep if you haven't already done so:    
  `go get github.com/tools/github`  
* Run Godep from the repository folder to vendor its dependencies:  
  `godep save`    
* Push to Predix Cloud Foundry  
  `cf push <your app name>`  

### Notes
#### Buildpacks
The default Go buildpack currently deployed in Predix is very out of date. Specify a custom buildpack when you push your app.
There are two ways to do this: 
* Add it to your manifest.yml:  
  `buildpack : https://github.com/cloudfoundry/go-buildpack.git`
* Use the `-b` option when you push your app:  
  `cf push <your-app-name> -b https://github.com/cloudfoundry/go-buildpack.git` 

#### Godep
The Go buildpack relies on Godep for dependency vendoring. Remember to update your dependencies before pushing to Predix. 
* `godep save` adds new dependencies, but does not update existing ones. 
* `godep update ...` updates all existing dependencies to the current version on your GOPATH.
* Don't use the `godep save -r` command to update your imports to use the vendored libraries. , The buildpack compile script runs `godep go install` which makes it unnecessary.

#### Target executable name and the run command
The built executable name is taken from the ImportPath in Godeps/Godeps.json, which defaults to the name of the source directory where main resides. That name must be specified exactly in the run command in your mainfest.yml (or with the `cf push -c` option). If your app uploads and compiles, but fails to run, check this first.  
Tip: Use the absolute path (`/app/bin/<your-app-name>`) to ensure you run your app, and not a different command with the same name appearing earlier in the path.      
    
    
## Working with UAA
These instructions assume you are developing a resource service. The requirements for a client service are somewhat different, and not currently covered by uaa-support or these notes.  
To get requests secured by UAA to succeed, you will need a UAA instance bound to this app with the following items defined/configured as a minimum:
* An audience for your resource server. This is not something to configure in UAA, but is used as a root for your scopes. In UAA, the audiences in an access token are auto-generated from the scopes by removing the final period and anything following it. (So in this example app a scope of `go-rest-service.read` has a derived audience of `go-rest-service`)
* SCIM groups for you resource server scopes. 
* A client account for your resource service, with "uaa.resource" authority to allow it to check token validity.
* For a client service making requests with user authorization:
  * A user account that is a member of the resource server SCIM groups.
  * A client account for your client service whose **scopes** include those SCIM groups.
* For a client service making requests under its own authorization:
  * A client account for your client service whose **authorities** include those SCIM groups.
  
You will also need to set UAA_CLIENT_ID and UAA_CLIENT_SECRET environment variables in your app to provide the resource service account credentials.

#### Testing
To make a UAA authorized request, you must include an Authorization header containing a bearer token in addition to any headers and body required by the request itself.

## Contributing
This code is intended to be a useful bootstrap for new services. If you find any bugs, or have suggestions for improvements or new features please let us know or (better) submit a pull request.
 

## Useful links
[Go official site](https://golang.org/)  
[The Go Tour](https://tour.golang.org/)  
[How to write Go](https://golang.org/doc/code.html)  
[Writing effective Go](https://golang.org/doc/effective_go.html)   
[Go official language spec](https://golang.org/ref/spec)

## Future Enhancements
* Swagger documentation
* AMQP interfaces
* RiakCS/S3 support access
