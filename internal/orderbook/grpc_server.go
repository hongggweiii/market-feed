package orderbook

import (
	"context"

	orderbookpb "github.com/hongggweiii/market-nodes/api/proto"
)

type GrpcServer struct {
	engine                                          *OrderBook
	orderbookpb.UnimplementedOrderBookServiceServer // Allow code to compile even if we don't implement all methods
}

func NewGrpcServer(engine *OrderBook) *GrpcServer {
	return &GrpcServer{engine: engine}
}

func (s *GrpcServer) GetOrderBook(ctx context.Context, req *orderbookpb.GetTopBookRequest) (*orderbookpb.GetTopBookResponse, error) {
	bestBidPrice, bestBidQty, bestAskPrice, bestAskQty := s.engine.GetTopBook()

	return &orderbookpb.GetTopBookResponse{
		Symbol:       req.GetSymbol(),
		BestBidPrice: bestBidPrice.String(),
		BestBidQty:   bestBidQty.String(),
		BestAskPrice: bestAskPrice.String(),
		BestAskQty:   bestAskQty.String(),
	}, nil
}
