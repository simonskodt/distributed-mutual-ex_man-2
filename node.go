package main

import (
	"context"
    "log"
    "os"
    "net"
	
    service "distributed-mutual-exclusion_mand-2/service"
    "google.golang.org/grpc"
)

type Connection struct {
    clientConn *grpc.ClientConn
    client     service.ProjectServiceClient
    context    context.Context
    port       string
}

type Node struct {
    nodeListningPort string
    clientList          map[string]service.ProjectServiceClient
    service.UnimplementedProjectServiceServer
}

var serverPorts = []string{} // args
var connections []Connection // serverConnections

//          Own
//          Port
//            |
//            V
// go run . 9000 9001 9002
// go run . 9001 9000 9002
// go run . 9002 9000 9001
func main() {
    userInput := os.Args[1:]
    ownPort := ":" + userInput[0]
    peer1 := ":" + userInput[1]
    peer2 := ":" + userInput[2]

    setupServer(ownPort)

    serverPorts = append(serverPorts, peer1, peer2)

    c := Node{nodeListningPort: ownPort}

    for i := range serverPorts {
        ctx, conn, c := setupConnection(i, &c)
        newConn := Connection{
            context:    ctx,
            clientConn: conn,
            client:     c,
            port:       serverPorts[i],
        }
        connections = append(connections, newConn)
        defer newConn.clientConn.Close()
    }

    for {
    }

}

func setupServer(port string) {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    server := grpc.NewServer()
    s := Node{}

    service.RegisterProjectServiceServer(server, &s)
    log.Printf("Server Port: %v \n", lis.Addr())
    go func() {
        err = server.Serve(lis)
        if  err != nil {
            log.Fatalf("Failed to serve at: %v", err)
        }
    }()
}

func setupConnection(index int, n *Node) (context.Context, *grpc.ClientConn, service.ProjectServiceClient) {
    context := context.Background()
    log.Printf("NODE LISTENING PORT: %v", n.nodeListningPort)

    conn, err := grpc.Dial(serverPorts[index], grpc.WithInsecure())

    if err != nil {
        log.Printf("Error: %v", err)
    }

    client := service.NewProjectServiceClient(conn)

    log.Printf("Connecting to: %v \n",serverPorts[index] )
    return context, conn, client
}