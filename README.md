# Sequra Backend Coding Challenge

This challenge was written in Google Go programming language. There are a few items I have not
addressed due to time constraints which I will outline below at the end of this README.

## Solution

| Year | Number of Disbursements | Amt Disbursed to Merchants | Amt of Order Fees | Number of Monthly Fees Charged | Amt of Monthly Fees Charged |
|------|-------------------------|----------------------------|-------------------|--------------------------------|-----------------------------|
| 2022 | 1,547                   | €36,433,527.69             | €1,419,169.02     | 29                             | €750.00                     |
| 2023 | 10,365                  | €182,392,281.43            | €7,290,126.22     | 148                            | €3,915.00                   |


## Setup and Run

The solution includes a sqlite3 database file which should be able to be run MacOS without any need for installation as MacOS ships with Sqlite3. If you 
are using Linux you may have to install sqlite3 using your distributions package manager. The sqlite3 DB schema is already configured, but I have included
the required createSqlTables.sql should you wish to recreate the tables. 

To run the service you can choose to open the project files in your preferred IDE and the main func in `main.go` within the project root, or you can build 
the binary by running `go build main.go` from the project root and then running the resulting binary. 

There are two API endpoints. One which can be triggered with an `HTTP GET` to `http://localhost:8080/import` which will parse the two provided csv files and 
insert the parsed data into the DISBURSEMENTS table. The other is to retrieve the requested report data and takes an `HTTP POST` to `http://localhost:8080/disbusrement` .
The post body MUST be in the form of a JSON object with the valid years for the report data. Below is an example.

`{
"Name": "Disbusement Report",
"YYYY": "2023"
}`

**NOTE** 
1. The importation process takes about 15 minutes to insert the disbursement records into the database. Until the process is complete, the disbursement report will be incorrect. 
2. The merchants.csv and orders.csv files will not be included in the submission, but must be present in the project root when run. 

## Assumptions and Tradeoffs 

1. The solution was built without any third party libraries. Only the Go standard lib was used with the assumption being 
that the security of the app would be of utmost importance and reliance on third party libraries can introduce unknown CVEs.
2. Sqlite3 was used for simplicity and testing. It should not be used for production and its performance is poor as 
compared to MariaDB, postgres, or any other "real" RDBMS.
3. The repository pattern was used to make it a relatively easy switch to another RDBMS. 
4. Time constraints caused me to have to stop development on this project. The request for a full production ready service 
in 6 hours of development time is not quite feasible. If I had more time I would have build more test cases and dockerized the solution for ease of deployment. 
5. I would have also added the new order processing via the API had I more time. That said, if you decide to move forward we can add that during the coding session.
6. The project needs more refactoring and code cleanup. The time constraints prevented further development. 
