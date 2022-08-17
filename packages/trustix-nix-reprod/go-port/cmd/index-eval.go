// Copyright (C) 2022 Tweag IO
//
// This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/trustix-nix-reprod/eval"
)

var indexEvalCommand = &cobra.Command{
	Use:   "index-eval",
	Short: "Index evaluation",
	RunE: func(cmd *cobra.Command, args []string) error {

		evalConfig := eval.NewConfig()
		evalConfig.Expr = "./pkgs.nix"

		ctx := context.Background()

		results, err := eval.Eval(ctx, evalConfig)
		if err != nil {
			panic(err)
		}

		for wrappedResult := range results {
			result, err := wrappedResult.Unwrap()
			if err != nil {
				panic(err)
			}

			fmt.Println(result)
		}

		return nil
	},
}
