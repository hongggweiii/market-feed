package main

import (
	"log"
	"net"
	"time"

	orderbookpb "github.com/hongggweiii/market-feed/api/proto"
	"github.com/hongggweiii/market-feed/internal/domain"
	"github.com/hongggweiii/market-feed/internal/exchange"
	"github.com/hongggweiii/market-feed/internal/orderbook"
	"google.golang.org/grpc"
)

func RunOrderBook(api orderbook.DepthProvider, symbol string, port string) {
	engine := orderbook.NewOrderBook()

	// Channel to receive depth updates from Websocket
	updates := make(chan *domain.DepthUpdate, 1000)

	go func() {
		err := api.StreamOrderBookDepthUpdates(symbol, updates)
		if err != nil {
			log.Fatalf("Stream stopped: %v", err)
		}
	}()

	log.Println("Websocket started...")

	snapshot, err := api.FetchDepthSnapshot(symbol)
	if err != nil {
		log.Fatalf("Failed to fetch order book: %v", err)
	}

	engine.Seed(snapshot)
	log.Printf("Engine seeded! Bids: %d, Asks: %d", len(engine.GetBids()), len(engine.GetAsks()))

	// Start gRPC server to serve order book data to clients
	go func() {
		listener, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("[%s] Failed to listen on port %s: %v", api.Name(), port, err)
		}
		defer listener.Close()

		grpcServer := grpc.NewServer()
		myServer := orderbook.NewGrpcServer(engine)

		// Regster the gRPC server
		orderbookpb.RegisterOrderBookServiceServer(grpcServer, myServer)

		log.Printf("[%s] gRPC server started on port %s", api.Name(), port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("[%s] Failed to serve gRPC: %v", api.Name(), err)
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case update := <-updates: // Channel receives new data from Websocket
			err := engine.ProcessUpdate(update)
			if err != nil {
				log.Fatalf("Fatal sync error: %v", err)
			}
		case <-ticker.C: // 1 second passed
			bestBidPrice, bestBidQty, bestAskPrice, bestAskQty := engine.GetTopBook()
			log.Printf("[%s] Best Bid: %s (%s), Best Ask: %s (%s)", api.Name(), bestBidPrice, bestBidQty, bestAskPrice, bestAskQty)
		}
	}
}

func main() {
	binanceAPI := &exchange.BinanceClient{}
	bybitAPI := &exchange.ByBitClient{}

	go RunOrderBook(binanceAPI, "BTCUSDT", ":50001")
	go RunOrderBook(bybitAPI, "BTCUSDT", ":50002")

	select {}
}
