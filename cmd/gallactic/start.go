package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gallactic/gallactic/core"
	"github.com/gallactic/gallactic/core/config"
	"github.com/gallactic/gallactic/core/proposal"
	"github.com/gallactic/gallactic/crypto"
	"github.com/gallactic/gallactic/keystore/key"
	"github.com/gallactic/gallactic/version"
	"github.com/jawher/mow.cli"
)

var welcomeMessage = `
       ___           ___                                       ___           ___                                   ___
      /  /\         /  /\                                     /  /\         /  /\          ___       ___          /  /\
     /  /:/_       /  /::\                                   /  /::\       /  /:/         /  /\     /  /\        /  /:/
    /  /:/ /\     /  /:/\:\    ___     ___   ___     ___    /  /:/\:\     /  /:/         /  /:/    /  /:/       /  /:/
   /  /:/_/::\   /  /:/~/::\  /__/\   /  /\ /__/\   /  /\  /  /:/~/::\   /  /:/  ___    /  /:/    /__/::\      /  /:/  ___
  /__/:/__\/\:\ /__/:/ /:/\:\ \  \:\ /  /:/ \  \:\ /  /:/ /__/:/ /:/\:\ /__/:/  /  /\  /  /::\    \__\/\:\__  /__/:/  /  /\
  \  \:\ /~~/:/ \  \:\/:/__\/  \  \:\  /:/   \  \:\  /:/  \  \:\/:/__\/ \  \:\ /  /:/ /__/:/\:\      \  \:\/\ \  \:\ /  /:/
   \  \:\  /:/   \  \::/        \  \:\/:/     \  \:\/:/    \  \::/       \  \:\  /:/  \__\/  \:\      \__\::/  \  \:\  /:/
    \  \:\/:/     \  \:\         \  \::/       \  \::/      \  \:\        \  \:\/:/        \  \:\     /__/:/    \  \:\/:/
     \  \::/       \  \:\         \__\/         \__\/        \  \:\        \  \::/          \__\/     \__\/      \  \::/
      \__\/         \__\/                                     \__\/         \__\/                                 \__\/    `

func Start() func(cmd *cli.Cmd) {
	return func(cmd *cli.Cmd) {

		workingDirOpt := cmd.String(cli.StringOpt{
			Name: "w working-dir",
			Desc: "working directory of the configuration files",
		})

		privatekeyOpt := cmd.String(cli.StringOpt{
			Name: "p privatekey",
			Desc: "private key of the node's validator",
		})

		keystoreOpt := cmd.String(cli.StringOpt{
			Name: "k key-file",
			Desc: "path to the encrypted node's key file",
		})

		keyfileauthOpt := cmd.String(cli.StringOpt{
			Name: "a auth",
			Desc: "key file passphrase",
		})

		/*
			cmd.Spec = "--working-dir=<working directory of the configuration files>"
		*/

		cmd.Action = func() {
			fmt.Print(welcomeMessage)
			fmt.Println("\n\n\nYou are running a gallactic blockchian node version: ", version.Version, ". Welcome!")
			workingDir := *workingDirOpt
			if workingDir != "" {
				keyObj := new(key.Key)
				switch {
				case *keystoreOpt == "" && *privatekeyOpt == "":
					// Creating KeyObject from Private Key
					kj, err := PromptPrivateKey()
					if err != nil {
						log.Fatalf("Aborted: %v", err)
					}
					keyObj = kj
				case *keystoreOpt != "" && *keyfileauthOpt != "":
					//Creating KeyObject from keystore
					passphrase := *keyfileauthOpt
					kj, err := key.DecryptKeyFile(*keystoreOpt, passphrase)
					if err != nil {
						log.Fatalf("Could not decrypt file: %v", err)
					}
					keyObj = kj
				case *keystoreOpt != "" && *keyfileauthOpt == "":
					//Creating KeyObject from keystore
					passphrase := promptPassphrase(true)
					kj, err := key.DecryptKeyFile(*keystoreOpt, passphrase)
					if err != nil {
						log.Fatalf("Could not decrypt file: %v", err)
					}
					keyObj = kj
				case *privatekeyOpt != "":
					// Creating KeyObject from Private Key
					pv, err := crypto.PrivateKeyFromString(*privatekeyOpt)
					if err != nil {
						log.Fatalf("Could not decrypt file: %v", err)
					}
					kj := CreateKey(pv)
					keyObj = kj
				}

				fmt.Println("Validator address: ", keyObj.Address().String())

				// change working directory
				if err := os.Chdir(workingDir); err != nil {
					log.Fatalf("Unable to changes working directory: %v", err)
				}
				configFile := "./config.toml"
				genesisFile := "./genesis.json"

				gen, err := proposal.LoadFromFile(genesisFile)
				if err != nil {
					log.Fatalf("Could not obtain genesis from file: %v", err)
				}

				conf, err := config.LoadFromFile(configFile)
				if err != nil {
					log.Fatalf("Could not obtain config from file: %v", err)
				}

				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()

				signer := crypto.NewValidatorSigner(keyObj.PrivateKey())
				kernel, err := core.NewKernel(ctx, gen, conf, signer)
				if err != nil {
					log.Fatalf("Could not create kernel: %v", err)
				}

				err = kernel.Boot()
				if err != nil {
					log.Fatalf("Could not boot kernel: %v", err)
				}

				kernel.WaitForShutdown()

			} else {
				fmt.Println("see 'gallactic start --help ' list of available commands to start gallactic node")
			}
		}
	}
}

func CreateKey(pv crypto.PrivateKey) *key.Key {
	addr := pv.PublicKey().ValidatorAddress()
	return key.NewKey(addr, pv)
}
