# pixicoreAPI


## How to run the API
- Install go

- Install all the package dependencies with  `go get -d ./...` inside project directory

- `go build && ./pixicoreAPI`

#### Using Docker

- `docker build -t pixicoreapi . `

- `docker run -d -p 3000:3000 pixicoreapi `

## Usage

- run `curl -i http://localhost:3000/v1/install/SERVER_MAC_ADDRESS`  this will collect info  and install coreOS for the server

## API Endpoints

#### `GET v1/boot/:macAddress` 

- Used by pixicore to get PXE config and boot each server (each server have a IP address assigned).

#### `GET v1/install/:macAddress` (NOT DONE) 

- get information (cores,ram,etc) from the server using her macAddress as ID and install coresOS. 

#### `GEt v1/install/all` (NOT DONE) *to work you have to change the IP addr of the servers in config file

- get information (cores,ram,etc) from each the server using her macAddress as ID and install coresOS for each one.

#### `GEt v1/info`

- show information about all the registred servers in DB.

### TO DO

- Unit tests
    - https://semaphoreci.com/community/tutorials/test-driven-development-of-go-web-applications-with-gin
    - https://medium.com/@craigchilds94/testing-gin-json-responses-1f258ce3b0b1