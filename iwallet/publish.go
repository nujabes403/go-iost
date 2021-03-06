// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package iwallet

/*
// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "sign a .sc file with .sig files, and publish it",
	Long:  `sign a .sc file with .sig files, and publish it`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(`invalid input, check
	iwallet publish -h`)
			return
		}

		sc, err := readFile(args[0])
		if err != nil {
			fmt.Println("Read file failed: ", err.Error())
			return
		}

		var mtx tx.Tx
		err = mtx.Decode(sc)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		signs := make([]*crypto.Signature, 0)
		for i, v := range args {
			if i == 0 {
				continue
			}
			sig, err := readFile(v)
			if err != nil {
				fmt.Println("Read file failed: ", err.Error())
				return
			}
			var sign crypto.Signature
			err = sign.Decode(sig)
			if err != nil {
				fmt.Println("Error: Illegal sig file", err)
				return
			}
			if !mtx.VerifySigner(&sign) {
				fmt.Printf("Error: Sign %v wrong\n", v)
				return
			}
			signs = append(signs, &sign)
		}
		fsk, err := readFile(kpPath)
		if err != nil {
			fmt.Println("Read file failed: ", err.Error())
			return
		}

		acc, err := account.NewKeyPair(loadBytes(string(fsk)), getSignAlgo(signAlgo))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		stx, err := tx.SignTx(&mtx, acc.ID, []*account.KeyPair{acc}, signs...)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		dest = changeSuffix(args[0], ".tx")

		saveTo(dest, stx.Encode())

		if !isLocal {
			var txHash string
			txHash, err = sendTx(stx)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("iost node:receive your tx!")
			fmt.Println("the transaction hash is:", txHash)
			if checkResult {
				checkTransaction(txHash)
			}
		}
	},
}


func init() {
	rootCmd.AddCommand(publishCmd)

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	publishCmd.Flags().StringVarP(&kpPath, "key-path", "k", home+"/.iwallet/id_ed25519", "Set path of sec-key")
	publishCmd.Flags().BoolVar(&isLocal, "local", false, "Set to not send tx to server")
	publishCmd.Flags().StringVarP(&signAlgo, "signAlgo", "a", "ed25519", "Sign algorithm")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// SetServer set the server.
func SetServer(s string) {
	server = s
}
*/
