# HTTP Protocol in Go

A from-scratch implementation of an HTTP/1.1 client written in Go.  
This project is part of learning Go and understanding how the HTTP protocol works beneath the standard library.

## Overview

This client handles the core behavior of HTTP/1.1 while relying on Go’s TCP facilities for the underlying connection. The project focuses on the protocol itself rather than reinventing TCP.

The client performs:

- Opening a TCP connection with Go’s `net` package
- Manually constructing and sending HTTP/1.1 request messages
- Supporting all HTTP methods (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD)
- Parsing the status line, headers, and body
- Following HTTP/1.1 formatting rules for requests and responses
