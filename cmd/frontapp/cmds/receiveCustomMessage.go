package cmds

import (
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/fairjungle/frontapp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var receiveCustomMessageCmd = &cobra.Command{
	Use:     "receive_custom_message",
	Short:   "Post a custom message in a custom channel",
	Example: `  frontapp receive_custom_message channel_id "Sender Handle" "This is the message body"`,
	RunE:    receiveCustomMessageRun,
}

func init() {
	receiveCustomMessageCmd.Flags().String("bodyFormat", "markdown", "Format of the message body")
	if err := viper.BindPFlag("bodyFormat", receiveCustomMessageCmd.Flags().Lookup("bodyFormat")); err != nil {
		panic(err)
	}

	receiveCustomMessageCmd.Flags().Bool("dump", false, "Dump parsed response")
	if err := viper.BindPFlag("dump", receiveCustomMessageCmd.Flags().Lookup("dump")); err != nil {
		panic(err)
	}

	receiveCustomMessageCmd.Flags().String("metaThreadRef", "", "Reference which will be used to thread messages. If omitted, Front threads by sender instead")
	if err := viper.BindPFlag("metaThreadRef", receiveCustomMessageCmd.Flags().Lookup("metaThreadRef")); err != nil {
		panic(err)
	}

	receiveCustomMessageCmd.Flags().String("senderContactId", "", "ID of the contact in Front corresponding to the sender")
	if err := viper.BindPFlag("senderContactId", receiveCustomMessageCmd.Flags().Lookup("senderContactId")); err != nil {
		panic(err)
	}

	receiveCustomMessageCmd.Flags().String("senderName", "", "Name of the sender")
	if err := viper.BindPFlag("senderName", receiveCustomMessageCmd.Flags().Lookup("senderName")); err != nil {
		panic(err)
	}

	receiveCustomMessageCmd.Flags().String("subject", "", "Subject of the message")
	if err := viper.BindPFlag("subject", receiveCustomMessageCmd.Flags().Lookup("subject")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(receiveCustomMessageCmd)
}

func receiveCustomMessageRun(cmd *cobra.Command, args []string) error {
	// init logger
	level, err := logrus.ParseLevel(viper.GetString("logLevel"))
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	log := logrus.New()
	log.Level = level

	// init client
	apiToken := viper.GetString("apiToken")
	if apiToken == "" {
		return errors.New("missing api key")
	}

	client := frontapp.NewClient(apiToken, frontapp.Logger(log))

	req := &frontapp.ReceiveCustomMessageReq{
		ChannelID: args[0],
		Sender: frontapp.MessageSender{
			ContactID: viper.GetString("senderContactId"),
			Name:      viper.GetString("senderName"),
			Handle:    args[1],
		},
		Subject:    viper.GetString("subject"),
		Body:       args[2],
		BodyFormat: viper.GetString("bodyFormat"),
	}

	if ref := viper.GetString("metaThreadRef"); ref != "" {
		req.Metadata = frontapp.MessageMetadata{
			ThreadRef: ref,
		}
	}

	resp, err := client.ReceiveCustomMessage(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if viper.GetBool("dump") {
		spew.Dump(resp)
	}

	return nil
}
