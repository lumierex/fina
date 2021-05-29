### frpc

#### Codec

1. read header
2. read body
3. write(*Header, body)

#### server
1. start server
2. read request
3. read header
4. read body
5. handle request -> send message to client 
