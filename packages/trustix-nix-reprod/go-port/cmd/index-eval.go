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
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tweag/trustix/packages/go-lib/set"
	drvparse "github.com/tweag/trustix/packages/trustix-nix-reprod/derivation"
	"github.com/tweag/trustix/packages/trustix-nix-reprod/eval"
	_ "modernc.org/sqlite"
)

const sqlDialect = "sqlite"

// func indexDrv(attr string, drv *derivation.Derivation

var indexEvalCommand = &cobra.Command{
	Use:   "index-eval",
	Short: "Index evaluation",
	RunE: func(cmd *cobra.Command, args []string) error {

		evalConfig := eval.NewConfig()
		evalConfig.Expr = "./pkgs.nix"

		ctx := context.Background()

		commitSha := "c4c79f09a599717dfd57134cdd3c6e387a764f63"

		fmt.Println(commitSha)

		db, err := sql.Open(sqlDialect, "./foo.sqlite3")
		if err != nil {
			return err
		}

		err = migrate(db, sqlDialect)
		if err != nil {
			panic(err)
		}

		results, err := eval.Eval(ctx, evalConfig)
		if err != nil {
			panic(err)
		}

		drvParser, err := drvparse.NewCachedDrvParser()
		if err != nil {
			panic(err)
		}

		// // Map drv to it's direct references
		// refs := make(map[string][]string)

		for wrappedResult := range results {
			result, err := wrappedResult.Unwrap()
			if err != nil {
				panic(err)
			}

			drv, err := drvParser.ReadPath(result.DrvPath)
			if err != nil {
				panic(err)
			}

			// Direct dependencies
			refsDirect := set.NewSet[string]()
			for inputDrv, _ := range drv.InputDerivations {
				refsDirect.Add(inputDrv)
			}

			// All dependencies (recursive, flattened)
			refsAll := refsDirect.Copy()

			// for inputDrv, _ := range drv.InputDerivations {
			// }

			fmt.Println(refsAll.Values())
		}

		return nil
	},
}
