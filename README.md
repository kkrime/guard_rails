# guard_rails

### Setup
The Database Schema files is in the `./setup` folder.<br/>
From the root folder for postgres you can run: `sudo -u postgres psql -c '\ir setup/schema.sql'`<br/>

### Database Schema
The database schema is very simple and it's located `./setup/schema.sql`

### Configuration
The configuation file is `./config.tml`, this file contains the main configuration file. Here you can configure the database and the repository scanners.
This file can easily be extended in the future to include things like the server port ect.

### To Start The Server
From the repository root run `go run .` and the server will run on port `8080`

#### To Configure The Database
<code>[postgres]
host     = 'localhost'
port     = 5432
user     = 'postgres'
password = 'postgres'
dbname   = 'guard_rails'</code>

### Quick Start
- Configure everything using the `./config` file
- Run `sudo -u postgres psql -c '\ir setup/schema.sql'` from the project root folder
- Start the server using `go run .`
- Add a repository: `curl -X POST localhost:8080/v1/repository -d '{"name": "testRepo", "url": "https://github.com/rtyley/small-test-repo"}'`
- Queue a scan: `curl -X POST localhost:8080/v1/repository/scan/testRepo`<br/>
**Please note:** if there's a crash during the cloding proccess, please manually delete the repository from the `./repositories` folder, before attempting to scan that repository again.
I ran out of time before I could fix this.
- Check scan results: `curl -X GET localhost:8080/v1/repository/scan/testRepo`

## APIs
### Repository
To Create: <br />
`POST /v1/repository` Body: `{"name": <repository_name> (alphanumeric only), "url": <url_address> (must include protocal prefix)}`<br />
To Read: <br />
`GET /v1/repository/<repository_name>`<br />
To Update (only the url can be updated):<br /> 
`PUT /v1/repository/<repository_name>` Body: `{"name": <repository_name>, "url": <url_address> (must include protocal prefix)}`<br />
To Delete:<br />
`Delete /v1/repository/<repository_name>`<br />

### Scan
To queue a scan:<br />
`POST /v1/repository/scan/<repository_name>`<br />
To get scan results (will get a list of all the scan results from most reccent):<br />
`GET /v1/repository/scan/<repository_name>`<br />

## Design
The code is split into 3 layers:<br />
**Controller Layer:** Responsible for parameter checks, data transformation and passing data to and from the user to the service layer<br />
**Service Layer:** This is where the business logic is, this is where the business logic lives<br />
**Database Layer:** Contains simple wrapper functions to write/read data to/from the database, minimal logic<br />

**Reason:** It helps to keep the code in seperate layers for testing and having each layer have specific responsibilites helps when issues/bugs need investigating.

### Scanning A Repository

To scan a repository, I created the `RepositoryScanner interface:`<br/><br/>
<code>type RepositoryScanner interface {
Scan(file client.File) *ScanResult
}
</code>

**Reason:** Any new additioanl repository scanners (of any type) can be added by implementing the interface above and configured using `./config` 

### Task
For the task, the repository needs to be scanned to look for text with the prefix `private_key` and `public_key`.<br/>
I decided to implement this using `regex`, which I called a `token` in the code.

#### TokenScanner
The `tokenScanner` is a `RepositoryScanner` that is configurable via `./config.toml`

<code>[[tokenScanner]]
[tokenScanner.scanData]
token = 'private_key\w+'
type = 'private key leak'
ruleId = '0001'
[tokenScanner.metadata]
description = 'private key leak'
severity = 'HIGH'
</code>

As you can see from above `token = 'private_key\w+'` is the `regex` that will be used to scan the repository to look for the `private_key` prefix
There are 2 `TokenScanners` configured, one for each prefix.<br/><br/>

Each file is scanned in its own `goroutine`.

The `TokenScanner` reads each file **one line at a time** <br/>
**Reason:** this is so that we do not flood the memory and cause a system crash.

### Clients
I created some wrapper clients interfaces in the `./client` folder. <br/>
**Reason:** I implemented these client interfaces is to help with unit testing

### Queue
If this was a application in proudction, we would use a `redis` or `kafka` queue, but due to the scrop of the task, the queue is implemented using a channel in go.
<br/>
When a scan request is made via `POST /v1/repository/scan/<repository_name>` the `scan` is added to the queue.<br/><br/>
The queue size is configured in the `./config` file:<br/>
<code>[queue]
queueSize = 4098</code>

### Design decisions
- I used integers for the table ids to keep things simple, in a real world application, I may have used `UUID`
- I didn't use a ORM for the database, reason being I personally find ORMs limiting and I prefer explicit SQL, it makes debugging easier for me

### Design Limitations
These are things I didn't have time to finish off, but in real world production would be needed:
- Startup check for scans `IN PROGRESS` and requeue them.
- If there's a crash duing the clone, subsequent scans will fail because the repository needs to be deleted and recloned. This has to be done manually for now.
- There is no way to tell if a scan failed due to a error or if the scan just failed organically. To rectify this a `errors` table should be added to the database which as a `repository_id` that references entries in the `repository` table. 
- There is no context for the scan implemented - using `context.Background()` for as a workaround.
- Need to add validation for `./config.toml`

## Testing

### End To End Tests

I have end to end tests in the `./e2e_tests/` folder. It's only one python script `repository_crud_e2e_test.py` - this is a crude script written for my sake more than anything else.
It was used as aprt of the development and is not part of a framework. In a real production environment this would be implemented via a framework and be part of the `CI/CD`.

### Unit Tests
I have unit tested the core functionality. Please find the unit tests:
- `./service/scan_service_test.go`
- `./scan/scan_test.go`

These tests aren't fully comprehensive, they are added to convey that I know how to write unit test and use mocks ect. In a real world prouction, the unit tests would provide full code coverage.

## Logging
Each log has an Id to help with debugging:
- All REST logs have a `RequestId` in the logs
- All Scan logs have the `ScanId` in the logs



