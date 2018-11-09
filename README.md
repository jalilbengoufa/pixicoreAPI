# pixicoreAPI

### API Endpoints

#### `GET v1/boot/:macAddress` 

- Used by pixicore to get PXE config and boot each server (each server have a IP address assigned).

#### `GET v1/install/:macAddress` (NOT DONE)

- get information (cores,ram,etc) from the server using her macAddress as ID and install coresOS. 

#### `GEt v1/install/all` (NOT DONE)

- get information (cores,ram,etc) from each the server using her macAddress as ID and install coresOS for each one.

#### `GEt v1/reset/:macAddress` (IN PROGRESS)

- remove information stored about a server from DB.

#### `GEt v1/reset/all`

- remove all information about servers from DB.

#### `GEt v1/ips`

- show all the static IPs available to use by servers.

#### `GEt v1/info`

- show information about all the registred servers in DB.