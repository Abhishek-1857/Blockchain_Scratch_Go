
# Superior Node
This is a superior node that receives transaction data from a central node and produces a block from that data. It then broadcasts the block to group of nodes.



## Features

- Receives transaction data from a central node

- Produces a block from the transaction data

- Broadcasts the block to groups of nodes

## Installation

To install the superior node, you will need to have the following dependencies installed:

-  Go : https://go.dev/doc/install

-  Fiber : https://docs.gofiber.io/ 

Once you have installed the dependencies, you can clone the repository and run the following command to build the superior node:

```go
git clone https://github.com/BlockmagixChain/blxtestnet.git

go run main.go
```


## Usage

To use the superior node, you will need to send it a POST request with the transaction data. The transaction data should be in JSON format.

The superior node will then produce a block from the transaction data and broadcast the block to groups of nodes.


## Configuration

You can configure the superior node by setting environment variables. The following is a list of environment variables that you can set:

- `PORT`: The port that the superior node will listen on

- `GROUPS`: The number of groups to divide the nodes into

- `PARAM`: The number of parameters to use in the group generation algorithm

## TO Do

- Develop a central node that collects transaction data  and sends it to the superior node.
- Implement validation code for groups of nodes to verify the integrity of the block before adding it to their respective blockchains.
- Explore integrating with different blockchains to expand the functionality of the superior node.
- Implement support for different transaction types to handle various transaction scenarios.
