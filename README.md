# brouker
Brouker is a broker between some form of output and a browser written in go.

It uses websockets 

The basic idea is that you post output to http://{host}:{port}/msg and it gets broadcased to all open connections, a la all the webchat examples you find online.
 
