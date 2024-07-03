package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	client := rpc.New(rpc.DevNet_RPC)
	wallet := solana.NewWallet()
	mint := solana.NewWallet()
	mintAuth := wallet.PublicKey()
	freezeAuth := mint.PublicKey()
	// get sol airdrop (24hr)
	out, err := client.RequestAirdrop(context.Background(), wallet.PublicKey(), solana.LAMPORTS_PER_SOL*2, rpc.CommitmentFinalized)
	if err != nil {
		log.Fatal("get airdrop err:", err)
	}
	fmt.Println(out)
	// init
	createInst, err := system.NewCreateAccountInstruction(0, 100, solana.TokenProgramID, mint.PublicKey(), wallet.PublicKey()).ValidateAndBuild()
	if err != nil {
		log.Fatal("create:", err)
	}
	// mint
	ins := token.NewInitializeMint2Instruction(9, mintAuth, freezeAuth, mint.PublicKey()).Build()
	// byt, _ := ins.Data()
	// tx, err := solana.TransactionFromBytes(byt)
	// if err != nil {
	// 	log.Fatal("ins ot tx:", err)
	// }
	recenthas, _ := client.GetRecentBlockhash(context.Background(), rpc.CommitmentFinalized)
	tx, err := solana.NewTransaction([]solana.Instruction{createInst, ins}, recenthas.Value.Blockhash, solana.TransactionPayer(wallet.PublicKey()))
	if err != nil {
		log.Fatal("new tx err:", err)
	}
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	},
	)
	if err != nil {
		log.Fatal(err)
	}

	sign, err := client.SendTransaction(context.TODO(), tx)
	if err != nil {
		log.Fatal("sign err:", err)
	}
	fmt.Println(sign)
}
