/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the cpackget project. */

package commands

import (
	"github.com/open-cmsis-pack/cpackget/cmd/cryptography"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var checksumCreateCmdFlags struct {
	// hashAlgorithm is the cryptographic hash function to be used
	hashAlgorithm string

	// outputDir is the target directory where the checksum file is written to
	outputDir string
}

var checksumVerifyCmdFlags struct {
	// checksumPath is the path of the checksum file
	checksumPath string
}

func init() {
	ChecksumCreateCmd.Flags().StringVarP(&checksumCreateCmdFlags.hashAlgorithm, "hash-function", "a", cryptography.Hashes[0], "specifies the hash function to be used")
	ChecksumCreateCmd.Flags().StringVarP(&checksumCreateCmdFlags.outputDir, "output-dir", "o", "", "specifies output directory for the checksum file")
	ChecksumVerifyCmd.Flags().StringVarP(&checksumVerifyCmdFlags.checksumPath, "path", "p", "", "path of the checksum file")

	ChecksumCreateCmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		err := command.Flags().MarkHidden("pack-root")
		_ = command.Flags().MarkHidden("concurrent-downloads")
		_ = command.Flags().MarkHidden("timeout")
		log.Debug(err)
		command.Parent().HelpFunc()(command, strings)
	})
	ChecksumVerifyCmd.SetHelpFunc(ChecksumCreateCmd.HelpFunc())
}

var ChecksumCreateCmd = &cobra.Command{
	Use:   "checksum-create [<local .path pack>]",
	Short: "Generates a .checksum file containing the digests of a pack",
	Long: `
Creates a .checksum file of a local pack. This file contains the digests
of the contents of the pack. Example <Vendor.Pack.1.2.3.sha256.checksum> file:

  "6f95628e4e0824b0ff4a9f49dad1c3eb073b27c2dd84de3b985f0ef3405ca9ca Vendor.Pack.1.2.3.pdsc
  435fsdf..."

  The referenced pack must be in its original/compressed form (.pack), and be present locally:

  $ cpackget checksum-create Vendor.Pack.1.2.3.pack

The default Cryptographic Hash Function used is "` + cryptography.Hashes[0] + `". In the future other hash functions
might be supported. The used function will be prefixed to the ".checksum" extension.

By default the checksum file will be created in the same directory as the provided pack.`,
	Args:              cobra.ExactArgs(1),
	PersistentPreRunE: configureInstallerVerbose,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cryptography.GenerateChecksum(args[0], checksumCreateCmdFlags.outputDir, checksumCreateCmdFlags.hashAlgorithm)
	},
}

var ChecksumVerifyCmd = &cobra.Command{
	Use:   "checksum-verify [<local .path pack>]",
	Short: "Verifies the integrity of a pack using its .checksum file",
	Long: `
Verifies the contents of a pack, checking its integrity against its .checksum file (created
with "checksum-create"), present in the same directory:

  $ cpackget checksum-verify Vendor.Pack.1.2.3.pack

The used hash function is inferred from the checksum filename, and if any of the digests
computed doesn't match the one provided in the checksum file an error will be thrown.
If the .checksum file is in another directory, specify it with the -p/--path flag`,
	Args:              cobra.ExactArgs(1),
	PersistentPreRunE: configureInstallerVerbose,
	RunE: func(cmd *cobra.Command, args []string) error {
		if checksumVerifyCmdFlags.checksumPath != "" {
			return cryptography.VerifyChecksum(args[0], checksumVerifyCmdFlags.checksumPath)
		}
		return cryptography.VerifyChecksum(args[0], "")
	},
}
