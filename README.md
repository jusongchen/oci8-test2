# oci8-test2

Step to reproduce the issue:

1)build the program

2) set env variable GO_OCI8_CONNECT_STRING
	export GO_OCI8_CONNECT_STRING=scott/tiger@//hostname/oracle_service

3) run the programm and enter CTRL+C after the prompt


4) CTRL+C signal won't cause any issue before a DB operation.  However, it will crash the programm with message "llegal instruction: 4" on Mac OS X after a DB operation.



