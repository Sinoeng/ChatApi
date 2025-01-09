
# Big Title

## Instructions

1. Make sure both `docker` and `docker-compose` is installed on your system. Make sure the daemon running. 
2. Clone the git directory and open it in a terminal. 
3. To start normally, run `docker-compose up --build -d`  
3.1 To start with the test container as well, run `docker-compose --profile test up --build -d`  
3.2 To read the results of the tests run `docker logs -ft primary_service_test` and read the results.  
3.3 On linux. To see what was tested you can run `go tool cover -html=coverage/coverage.html` from the projects root directory. This requires having `go` installed.  
4. To shut it down, run `docker-compose down`.  
4.1 If you ran the tests as well, run `docker-compose --profile test down` instead.  
5. Connect via localhost:8080

There is a default admin user called bob_admin with password bob.
You can find the api routes in /swagger/index.html, all routes have /v1/ as a prefix.
