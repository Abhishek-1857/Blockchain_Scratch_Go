package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"pop_v1/accounts"
	"pop_v1/accounts/accounthelper"
	"pop_v1/database"
	Ecdsa "pop_v1/ecdsa"
	"pop_v1/models"
	"pop_v1/router"
	transactionhelper "pop_v1/transaction.controller/transaction.helper"
	"pop_v1/utils"
	"pop_v1/wallet"
	"strconv"
	"strings"
	"time"

	// mydisc "pop_v1/discovery"

	//"pop_v1/router"

	//"github.com/gofiber/fiber/v2"

	"github.com/libp2p/go-libp2p"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
	"github.com/tecbot/gorocksdb"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other  peers.
const DiscoveryServiceTag = "pubsub"

const TopicName = "BlockMagix"
const AddTransaction = "Transaction"
const AddAccount = "Account"
const brodcast = "blockchain"

var (
	topicNameFlag = flag.String("topicName", "BlockMagix", "name of topic to join")
	name          string
	client_id     string
	client_addr   string
)

func main() {
	// parse some flags to set our nickname and the room to join
	opts1 := gorocksdb.NewDefaultOptions()
	opts2 := gorocksdb.NewDefaultOptions()
	opts3 := gorocksdb.NewDefaultOptions()
	opts4 := gorocksdb.NewDefaultOptions()

	defer opts1.Destroy()
	defer opts2.Destroy()
	defer opts3.Destroy()
	defer opts4.Destroy()

	opts1.SetCreateIfMissing(true)
	opts2.SetCreateIfMissing(true)
	opts3.SetCreateIfMissing(true)
	opts4.SetCreateIfMissing(true)

	database.Account_db, _ = gorocksdb.OpenDb(opts1, "database/accounts")
	database.Blockchain_db, _ = gorocksdb.OpenDb(opts2, "database/blockchain-db")
	database.Mempool_db, _ = gorocksdb.OpenDb(opts3, "database/mempool")
	database.TestTransaction_db, _ = gorocksdb.OpenDb(opts4, "database/testtransaction")

	flag.StringVar(&name, "name", "", " name to distinguish nodes")
	flag.StringVar(&client_id, "client_id", "", "id of the client node")
	flag.StringVar(&client_addr, "client_addr", "", "multi-address of the client node")
	flag.Parse()
	if name == "" {
		log.Fatal("You need to specify '-name' of your node to either 'client' , 'superior' or 'validator'")
	}

	if name != "client" {
		if client_id == "" {
			log.Fatal("You need to specify '-client_id'")
		}

		if client_addr == "" {
			log.Fatal("You need to specify '-client-addr'")
		}
		addr, _ := multiaddr.NewMultiaddr(client_addr)
		peerId, _ := peer.Decode(client_id)
		utils.Client_id = peerId
		utils.Client_addr = addr
	} else if name != "superior" && name != "validator" && name != "client" && name != "user" {
		log.Fatal("You need to specify '-name' of your node to either 'client' , 'superior' or 'validator' or 'user'")
	}

	if name == "user" {

		// Create an HTTP server with a single endpoint
		http.HandleFunc("/getaccount", func(w http.ResponseWriter, r *http.Request) {
			// Print the URL to the console
			// Start client
			clientHost, _ := libp2p.New(libp2p.NoListenAddrs)

			addr := utils.Client_addr
			peerId := utils.Client_id
			info := peer.AddrInfo{
				ID:    peerId,
				Addrs: []multiaddr.Multiaddr{addr},
			}

			clientHost.Connect(context.Background(), info)

			tr := &http.Transport{}
			tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))
			client := &http.Client{Transport: tr}

			res, err := Ecdsa.GenerateECDSAKeyPair(client, info.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			clientHost.Close()
			w.Write([]byte(res))
		})
		http.HandleFunc("/getbalance", func(w http.ResponseWriter, r *http.Request) {
			address := r.URL.Query().Get("address")
			clientHost, _ := libp2p.New(libp2p.NoListenAddrs)

			addr := utils.Client_addr
			peerId := utils.Client_id
			info := peer.AddrInfo{
				ID:    peerId,
				Addrs: []multiaddr.Multiaddr{addr},
			}

			clientHost.Connect(context.Background(), info)

			tr := &http.Transport{}
			tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))
			client := &http.Client{Transport: tr}
			// Validate and process the address as needed
			if address == "" {
				http.Error(w, "Missing 'address' parameter", http.StatusBadRequest)
				return
			}

			res, err := accounthelper.Retrievebalance(address, client, info.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			clientHost.Close()
			w.Write([]byte(res))

		})
		http.HandleFunc("/gettransaction", func(w http.ResponseWriter, r *http.Request) {
			clientHost, _ := libp2p.New(libp2p.NoListenAddrs)

			addr := utils.Client_addr
			peerId := utils.Client_id
			info := peer.AddrInfo{
				ID:    peerId,
				Addrs: []multiaddr.Multiaddr{addr},
			}

			clientHost.Connect(context.Background(), info)

			tr := &http.Transport{}
			tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))
			client := &http.Client{Transport: tr}
			hash := r.URL.Query().Get("hash")

			// Validate and process the address as needed
			if hash == "" {
				http.Error(w, "Missing 'hash' parameter", http.StatusBadRequest)
				return
			}

			res, err := transactionhelper.RetrievetransactionByHash(hash, client, info.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			clientHost.Close()
			w.Write([]byte(res))

		})

		http.HandleFunc("/sendtransaction", func(w http.ResponseWriter, r *http.Request) {
			// Extract the address from the URL path
			// Get the value of the "address" query parameter
			clientHost, _ := libp2p.New(libp2p.NoListenAddrs)

			addr := utils.Client_addr
			peerId := utils.Client_id
			info := peer.AddrInfo{
				ID:    peerId,
				Addrs: []multiaddr.Multiaddr{addr},
			}

			clientHost.Connect(context.Background(), info)

			tr := &http.Transport{}
			tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))
			client := &http.Client{Transport: tr}
			from := r.URL.Query().Get("from")

			to := r.URL.Query().Get("to")

			amount := r.URL.Query().Get("amount")

			// Validate and process the address as needed
			if from == "" {
				http.Error(w, "Missing 'from' parameter", http.StatusBadRequest)
				return
			}
			if to == "" {
				http.Error(w, "Missing 'to' parameter", http.StatusBadRequest)
				return
			}

			if amount == "" {
				http.Error(w, "Missing 'amount' parameter", http.StatusBadRequest)
				return
			}

			f, err := strconv.ParseFloat(amount, 64)
			if err != nil {
				log.Fatal("Error:", err)
				return
			}
			res, err := transactionhelper.SendTransaction(from, to, f, client, info.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			clientHost.Close()
			w.Write([]byte(res))
		})

		// Start the HTTP server
		log.Fatal(http.ListenAndServe(":8080", nil))

	} else {
		ctx := context.Background()
		// create a new libp2p Host that listens on a random TCP port
		serverHost, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
		if err != nil {
			panic(err)
		}

		log.Printf("Host ID: %s", serverHost.ID())
		log.Printf("Connect to me on:")
		for _, addr := range serverHost.Addrs() {
			log.Printf("  %s/p2p/%s", addr, serverHost.ID())
			if !strings.Contains(addr.String(), "127.0.0.1") {
				utils.Self_addr = addr
			}
		}

		fmt.Print(utils.Self_addr)

		// Initialize self host id
		utils.Self_id = serverHost.ID()

		//setup local mDNS discovery
		if err := setupDiscovery(serverHost); err != nil {
			panic(err)
		}

		// create a new PubSub service using the GossipSub router
		ps, err := pubsub.NewGossipSub(ctx, serverHost)
		if err != nil {
			panic(err)
		}
		utils.Ps = ps
		utils.Serverhost = serverHost
		// join the room from the cli flag, or the flag default
		utils.Topic, err = ps.Join(TopicName)
		if err != nil {
			panic(err)
		}
		utils.CTopic, err = ps.Join(AddTransaction)
		if err != nil {
			panic(err)
		}
		utils.AccountTopic, err = ps.Join(AddAccount)
		if err != nil {
			panic(err)
		}
		utils.BlockchainTopic, err = ps.Join(brodcast)
		if err != nil {
			panic(err)
		}

		listener, _ := gostream.Listen(serverHost, p2phttp.DefaultP2PProtocol)
		defer listener.Close()
		go func() {
			router.MainRoute()
			server := &http.Server{}
			server.Serve(listener)
		}()

		sub, err := utils.Topic.Subscribe()

		if err != nil {
			panic(err)
		}

		subC, errC := utils.CTopic.Subscribe()

		if errC != nil {
			panic(err)
		}

		subAccount, errAccount := utils.AccountTopic.Subscribe()

		if errAccount != nil {
			panic(err)
		}

		subBlockchain, errBlockchain := utils.BlockchainTopic.Subscribe()
		if errBlockchain != nil {
			panic(err)
		}
		time.Sleep(1 * time.Second)
		// showPeers(ps, serverHost)
		if name == "superior" {
			go BlockCreation()
		}
		go printMessagesFrom(ctx, sub)

		go addAccount(ctx, subAccount)

		go broadcastBlock(ctx, subBlockchain)

		addTransactions(ctx, subC)
	}
}

func BlockCreation() {
	time.Sleep(1 * time.Second)
	transactionhelper.TransactionLogic()
}

func broadcastBlock(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		var block models.Block

		err = json.Unmarshal([]byte(m.Message.Data), &block)
		if err != nil {
			log.Fatal(err)
		}
		blockId := utils.GenerateId()

		data, err := json.Marshal(block)
		if err != nil {
			log.Fatal(err)
		}

		// Writing data to the db
		writeOpts := gorocksdb.NewDefaultWriteOptions()
		err = database.Blockchain_db.Put(writeOpts, []byte(blockId), data)
		if err != nil {
			log.Fatal("Error writing data: in blockchain-db", err)
			return
		}

		log.Println("Block added successfully with blockheight : ", block.BlockHeader.Height, " blockhash : ", block.BlockHash)
		writeOpts.Destroy()
		opts_txn3 := gorocksdb.NewDefaultOptions()
		opts_txn3.SetCreateIfMissing(true)

		readOptions1 := gorocksdb.NewDefaultReadOptions()
		readOptions2 := gorocksdb.NewDefaultReadOptions()
		writeOption1 := gorocksdb.NewDefaultWriteOptions()
		writeOption2 := gorocksdb.NewDefaultWriteOptions()
		writeOption3 := gorocksdb.NewDefaultWriteOptions()

		for _, hash := range block.MetaData {
			// Read data from the database using the ID

			data1, err := database.TestTransaction_db.Get(readOptions1, []byte(hash))
			if err != nil {
				log.Fatalf("Error reading data from db:")
			}

			var trans models.Transaction
			if err := json.Unmarshal(data1.Data(), &trans); err != nil {
				log.Fatal("failed to unmarshal the transaction !!")
			}
			data1.Free()

			err = database.Mempool_db.Delete(writeOption2, []byte(hash))
			if err != nil {
				log.Fatal("failed to delete the transaction from the mempool !!")
			}

			trans.Txn.Status = "SUCCESS"
			transjson, err := json.Marshal(trans)
			if err != nil {
				log.Fatal("Error while Marshaling !!")
			}
			err = database.TestTransaction_db.Put(writeOption1, []byte(hash), transjson)
			if err != nil {
				log.Fatal("Error updating the transaction db !!")
			}

			from := trans.Txn.From
			to := trans.Txn.To
			var from_acc accounts.Account
			var to_acc accounts.Account

			data3, err := database.Account_db.Get(readOptions2, []byte(from))
			if err != nil {
				log.Fatalf("Error reading data from db:")
			}

			if data3.Size() > 0 {
				//Unmarshal the data into a Account struct
				if err := json.Unmarshal(data3.Data(), &from_acc); err != nil {
					log.Fatal("Account failed to unmarshal !!")
				}
			}
			data3.Free()

			data4, err := database.Account_db.Get(readOptions2, []byte(to))
			if err != nil {
				log.Fatal("Error reading data from db:")
			}
			if data4.Size() > 0 {
				//Unmarshal the data into a Account struct
				if err := json.Unmarshal(data4.Data(), &to_acc); err != nil {
					log.Fatal("Account failed to unmarshal !!")
				}
			}
			data4.Free()
			from_acc.Amount = from_acc.Amount - trans.Txn.Amount
			from_acc.Nonce += 1
			from_acc.Txn_pending = false

			to_acc.Amount = to_acc.Amount + trans.Txn.Amount
			fromjson, err := json.Marshal(from_acc)
			if err != nil {
				log.Fatalf("Error Marshaling data !!")
			}
			tojson, err := json.Marshal(to_acc)
			if err != nil {
				log.Fatalf("Error Marshaling data !!")
			}
			database.Account_db.Put(writeOption3, []byte(from), fromjson)
			database.Account_db.Put(writeOption3, []byte(to), tojson)
		}
		opts_txn3.Destroy()
		readOptions1.Destroy()
		readOptions2.Destroy()
		writeOption1.Destroy()
		writeOption2.Destroy()
		writeOption3.Destroy()
		if name == "superior" {
			go BlockCreation()
		}
	}
}

func addAccount(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		var acc_wallet wallet.Wallet

		err = json.Unmarshal([]byte(m.Message.Data), &acc_wallet)
		if err != nil {
			log.Fatal(err)
		}

		pubkey := acc_wallet.PublicKey

		var account accounts.Account

		account.Amount = 10000
		account.Denom = "blx"
		account.Nonce = 0
		account.PubKey = pubkey
		account.Txn_pending = false

		pub_addr := fmt.Sprintf("%s", acc_wallet.Address())
		accountJSON, err := json.Marshal(account)
		if err != nil {
			log.Fatal("Error serializing node", err)
		}
		// Writing data to the db
		writeOpts := gorocksdb.NewDefaultWriteOptions()

		err = database.Account_db.Put(writeOpts, []byte(pub_addr), accountJSON)
		if err != nil {
			log.Fatal("Error writing data: in accounts db", err)
			return
		}
	}
}

func showPeers(ps *pubsub.PubSub, node host.Host) {
	peers := ps.ListPeers(TopicName)
	log.Println("Connected peers : ")
	for _, p := range peers {
		log.Printf("%s", p)
		multiaddrs := node.Peerstore().Addrs(p)
		for _, addr := range multiaddrs {
			if !strings.Contains(addr.String(), "127.0.0.1") {
				log.Printf("%s", addr.String())
			}
		}
	}
}

func addTransactions(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		var transaction models.Transaction
		err = json.Unmarshal([]byte(m.Message.Data), &transaction)
		if err != nil {
			log.Fatal(err)
		}
		from_acc, _ := accounthelper.Getaccountbyaddress(transaction.Txn.From)
		from_acc.Txn_pending = true
		accountJSON, err := json.Marshal(from_acc)
		if err != nil {
			log.Fatal("Error serializing account", err)
		}
		// Writing data to the db
		writeOpts_acc := gorocksdb.NewDefaultWriteOptions()

		err = database.Account_db.Put(writeOpts_acc, []byte(transaction.Txn.From), accountJSON)
		if err != nil {
			log.Fatal("Error writing data: in accounts db", err)
		}

		transactionJSON, err := json.Marshal(transaction)
		if err != nil {
			log.Fatal("Error serializing node", err)
		}
		// Writing data to the db
		writeOpts_txn1 := gorocksdb.NewDefaultWriteOptions()
		writeOpts_txn2 := gorocksdb.NewDefaultWriteOptions()
		err = database.TestTransaction_db.Put(writeOpts_txn1, []byte(transaction.Txn.Txnhash), transactionJSON)
		if err != nil {
			log.Fatal("Error writing data: in testtransaction db", err)
		}

		err = database.Mempool_db.Put(writeOpts_txn2, []byte(transaction.Txn.Txnhash), transactionJSON)
		if err != nil {
			log.Fatal("Error writing data: in mempool db", err)
		}
		writeOpts_txn1.Destroy()
		writeOpts_txn2.Destroy()
		writeOpts_acc.Destroy()
	}
}

func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		log.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
	}
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	ps := utils.Ps
	// // node := utils.Serverhost
	peers := ps.ListPeers(TopicName)
	currentId := pi.ID.String()
	for _, p := range peers {
		// 	// multiaddrs := node.Peerstore().Addrs(p)
		// 	// for _, addr := range multiaddrs {
		// 	// 	if !strings.Contains(addr.String(), "127.0.0.1") {
		// 	// 		log.Printf("%s", addr.String())
		// 	// 	}
		// 	// }
		if p.String() == currentId {
			return
		}
	}
	log.Printf("discovered new peer %s\n", pi.ID)
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		log.Printf("error connecting to peer %s: %s\n", pi.ID, err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}

func ShowTransactions() {
	// Reading data
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	// Iterating through the database
	iter := database.TestTransaction_db.NewIterator(readOpts)
	defer iter.Close()

	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()
		client := models.Transaction{}
		if err := json.Unmarshal(value.Data(), &client); err != nil {
			log.Fatal("Error deserializing data", err)
		}
		log.Printf("%v : %v\n", string(key.Data()), client)
	}
}
