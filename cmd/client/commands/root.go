package commands

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/diSpector/currency.git/internal/currency/client"
	"github.com/diSpector/currency.git/pkg/currency/entities"
	pb "github.com/diSpector/currency.git/pkg/currency/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var rootCmd = &cobra.Command{
	Use:   "currency",
	Short: "get currency rate",
	Long:  `Get specified currencies rates; could take array divided by spaces`,
	Args:  ArgsValidator,
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("localhost:50060", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Println(`err dial grpc:`, err)
			os.Exit(1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client := pb.NewCurrencyApiClient(conn)

		stream, err := client.GetCurrency(ctx, &pb.CurrencyRequest{
			Name: args,
		})
		if err != nil {
			log.Println(`err get stream`, err)
			return
		}

		var curMap = make(map[string]*entities.Currency)
		for i := range args {
			curMap[args[i]] = nil
		}

		for {
			cur, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln("err read currency from stream:", err)
			}
			if _, ok := curMap[cur.Code]; !ok {
				log.Println(`unexpected currency code - `, cur.Code)
			} else {
				curMap[cur.Code] = &entities.Currency{
					Name: cur.Name,
					Code: cur.Code,
					Rate: cur.Rate,
				}
			}
		}

		for k, v := range curMap {
			if v != nil {
				log.Printf("currency - %s, rate = %f\n", v.Code, v.Rate)
			} else {
				log.Printf("currency %s NOT found \n", k)
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ArgsValidator(cmd *cobra.Command, args []string) error {
	var err error

	if err = cobra.RangeArgs(1, 10)(cmd, args); err != nil {
		return err
	}

	if err = client.ValidateCurrencyArgs(args); err != nil {
		return err
	}

	return nil
}
