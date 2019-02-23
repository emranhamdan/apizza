// Copyright © 2019 Harrison Brown harrybrown98@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"apizza/dawg"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Order pizza from dominos",
	RunE: func(cmd *cobra.Command, args []string) error {
		var order *dawg.Order

		cached, err := cmd.Flags().GetBool("cached")
		if err != nil {
			return err
		}

		if cached {
			fmt.Println("this is where you would see previous orders and saved orders")
			return nil
		}

		if order == nil {
			return errors.New("Error: no orders were given")
		}
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

var newOrderCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new order that will be stored in the cache.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if store == nil {
			store, err = dawg.NearestStore(addr, cfg.Service)
			if err != nil {
				return err
			}
		}
		order := store.NewOrder()

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		if name == "" {
			return errors.New("Error: No order name... use '--name=<order name>'")
		} else if name == "new" { // I completely forgot why I did this
			return errors.New("Error: cannot give an order that name")
		}

		product, err := cmd.Flags().GetString("product")
		if err != nil {
			return err
		}
		if product != "" {
			p, err := store.GetProduct(product)
			if err != nil {
				return err
			}
			order.AddProduct(p)
		}

		raw, err := json.Marshal(&order)
		if err != nil {
			return err
		}
		fmt.Print(name, ": ")
		fmt.Println(string(raw))
		print("\n\n")
		fmt.Printf("%+v\n", store)
		price, err := order.Price()
		if err != nil {
			return err
		}
		fmt.Println("Price:", price)

		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	newOrderCmd.Flags().StringP("name", "n", "", "set the name of a new order")
	newOrderCmd.Flags().StringP("product", "p", "", "the product code for the new order")

	orderCmd.AddCommand(newOrderCmd)

	orderCmd.Flags().BoolP("cached", "c", false, "show the previously cached and saved orders")
	rootCmd.AddCommand(orderCmd)
}