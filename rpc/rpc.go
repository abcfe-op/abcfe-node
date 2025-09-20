package rpc

import (
	"context"
	"fmt"
	"net"

	"github.com/abcfe-op/abcfe-node/blockchain"
	"github.com/abcfe-op/abcfe-node/common/utils"
	"github.com/abcfe-op/abcfe-node/proto"
	"github.com/abcfe-op/abcfe-node/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/abcfe-op/abcfe-node/common/logger"
)

type server struct {
	proto.UnimplementedBlockchainServiceServer
	port string
}

// func (s *server) GetBlocks(ctx context.Context, empty *proto.Empty) (*proto.BlocksResponse, error) {
// 	blocks := blockchain.Blocks(blockchain.Blockchain())
// 	protoBlocks := make([]*proto.Block, 0, len(blocks))

// 	// 블록체인 데이터를 proto 형식으로 변환
// 	for _, block := range blocks {
// 		protoBlock := &proto.Block{
// 			Hash:      block.Hash,
// 			Height:    int64(block.Height),
// 			Timestamp: block.Timestamp,
// 		}
// 		protoBlocks = append(protoBlocks, protoBlock)
// 	}

// 	return &proto.BlocksResponse{Blocks: protoBlocks}, nil
// }

func (s *server) GetBlocks(ctx context.Context, empty *proto.Empty) (*proto.BlocksResponse, error) {
	blocks := blockchain.Blocks(blockchain.Blockchain())
	protoBlocks := make([]*proto.Block, 0, len(blocks))

	for _, b := range blocks {
		protoBlock := &proto.Block{
			Hash:        b.Hash,
			PrevHash:    b.PrevHash,
			Height:      int32(b.Height),
			Timestamp:   int32(b.Timestamp),
			Transaction: setTransactions(b.Transaction),
			RoleInfo:    setProtoRoleInfo(b.RoleInfo),
			Signature:   setSignatures(b.Signature),
		}
		protoBlocks = append(protoBlocks, protoBlock)
	}

	return &proto.BlocksResponse{
		Blocks: protoBlocks,
	}, nil
}

func (s *server) GetBlock(ctx context.Context, req *proto.BlockRequest) (*proto.BlockResponse, error) {
	block, err := blockchain.FindBlock(req.Hash)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Block not found")
	}

	protoBlock := &proto.Block{
		Hash:        block.Hash,
		PrevHash:    block.PrevHash,
		Height:      int32(block.Height),
		Timestamp:   int32(block.Timestamp),
		Transaction: setTransactions(block.Transaction),
		RoleInfo:    setProtoRoleInfo(block.RoleInfo),
		Signature:   setSignatures(block.Signature),
	}

	return &proto.BlockResponse{
		Block: protoBlock,
	}, nil
}

func (s *server) GetStatus(ctx context.Context, req *proto.Empty) (*proto.StatusResponse, error) {
	blockchain.Blockchain()
	return &proto.StatusResponse{
		CurrentHeight: int64(blockchain.Blockchain().Height),
		CurrentHash:   blockchain.Blockchain().NewestHash,
	}, nil
}

func (s *server) GetWallet(ctx context.Context, req *proto.Empty) (*proto.WalletResponse, error) {
	wallet := wallet.Wallet(s.port[1:])
	return &proto.WalletResponse{
		Address: wallet.Address,
		Port:    s.port,
	}, nil
}

func (s *server) GetBalance(ctx context.Context, req *proto.BalanceRequest) (*proto.BalanceResponse, error) {
	amount := blockchain.BalanceByAddress(req.Address, blockchain.Blockchain())
	return &proto.BalanceResponse{
		Address: req.Address,
		Balance: int64(amount),
	}, nil
}

func (s *server) GetStakingList(ctx context.Context, req *proto.Empty) (*proto.StakingListResponse, error) {
	_, stakingWalletTx, _ := blockchain.UTxOutsByStakingAddress(utils.StakingAddress, blockchain.Blockchain())
	stakingInfoList := blockchain.GetStakingList(stakingWalletTx, blockchain.Blockchain())
	protoStakingList := make([]*proto.StakingInfo, 0, len(stakingInfoList))

	for _, s := range stakingInfoList {
		staker := &proto.StakingInfo{
			Hash:      s.ID,
			Address:   s.Address,
			Port:      s.Port,
			Timestamp: int32(s.TimeStamp),
		}
		protoStakingList = append(protoStakingList, staker)
	}
	return &proto.StakingListResponse{StakingList: protoStakingList}, nil
}

func (s *server) GetMempool(ctx context.Context, req *proto.Empty) (*proto.MempoolResponse, error) {
	txs := blockchain.Mempool().Txs
	protoTxs := make([]*proto.Transaction, 0, len(txs))
	for _, tx := range txs {
		protoTxs = append(protoTxs, setTransactions([]*blockchain.Tx{tx})[0])
	}
	return &proto.MempoolResponse{Transactions: protoTxs}, nil
}

// func (s *server) GetRoleInfo(ctx context.Context, req *proto.Empty) (*proto.RoleInfoResponse, error) {

// }

func Start(port int) {
	grpcPort := fmt.Sprintf(":%d", port+3333) // REST 포트 + 3333을 gRPC 포트로 사용
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Error(fmt.Sprintf("failed to listen: %v", err))
	}

	s := grpc.NewServer()
	proto.RegisterBlockchainServiceServer(s, &server{port: fmt.Sprintf(":%d", port)})

	fmt.Printf("gRPC Server listening on port %s\n", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Error(fmt.Sprintf("failed to serve: %v", err))
	}
}

func setProtoRoleInfo(info *blockchain.RoleInfo) *proto.RoleInfo {
	if info == nil {
		return nil
	}
	return &proto.RoleInfo{
		ProposerAddress:         info.ProposerAddress,
		ProposerPort:            info.ProposerPort,
		ProposerSelectedHeight:  int32(info.ProposerSelectedHeight),
		ValidatorAddress:        info.ValidatorAddress,
		ValidatorPort:           info.ValidatorPort,
		ValidatorSelectedHeight: int32(info.ValidatorSelectedHeight),
	}
}

func setSignatures(sigs []*blockchain.ValidateSignature) []*proto.ValidateSignature {
	if sigs == nil {
		return nil
	}
	protoSigs := make([]*proto.ValidateSignature, 0, len(sigs))
	for _, sig := range sigs {
		protoSig := &proto.ValidateSignature{
			Port:      sig.Port,
			Address:   sig.Address,
			Signature: sig.Signature,
		}
		protoSigs = append(protoSigs, protoSig)
	}
	return protoSigs
}

func setTransactions(txs []*blockchain.Tx) []*proto.Transaction {
	if txs == nil {
		return nil
	}
	protoTxs := make([]*proto.Transaction, 0, len(txs))
	for _, tx := range txs {
		protoTx := &proto.Transaction{
			Id:        tx.ID,
			Timestamp: int32(tx.Timestamp),
			TxIns:     setTxIns(tx.TxIns),
			TxOuts:    setTxOuts(tx.TxOuts),
			InputData: tx.InputData,
		}
		protoTxs = append(protoTxs, protoTx)
	}
	return protoTxs
}

func setTxIns(txIns []*blockchain.TxIn) []*proto.TxIn {
	if txIns == nil {
		return nil
	}
	protoTxIns := make([]*proto.TxIn, 0, len(txIns))
	for _, txIn := range txIns {
		protoTxIn := &proto.TxIn{
			TxId:      txIn.TxID,
			Index:     int32(txIn.Index),
			Signature: txIn.Signature,
		}
		protoTxIns = append(protoTxIns, protoTxIn)
	}
	return protoTxIns
}

func setTxOuts(txOuts []*blockchain.TxOut) []*proto.TxOut {
	if txOuts == nil {
		return nil
	}
	protoTxOuts := make([]*proto.TxOut, 0, len(txOuts))
	for _, txOut := range txOuts {
		protoTxOut := &proto.TxOut{
			Address: txOut.Address,
			Amount:  int32(txOut.Amount),
		}
		protoTxOuts = append(protoTxOuts, protoTxOut)
	}
	return protoTxOuts
}
