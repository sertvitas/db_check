- checks connectivity to a database instance
- tries at regular intervals to connect to the database (interval time?)
- login and dbping
- log success/fail with timestamp for each interval
- --starttime of overall process
- --outage time(s)
- --end time of overall process





```bash
export AWS_REGION=us-east-1
./db_check <secret ID>


```