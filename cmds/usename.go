package cmds

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/cli"
)

// type response struct {
// }

type Username struct{ UI cli.ColoredUi }

func validateInputArgs(args []string) error {
	if len(args) != 1 {
		err := errors.New("please enter valid argument")
		return err
	}
	return nil
}

// Help implements cli.Command.
func (b *Username) Help() string {
	desc := "Usage: github-activity username <username>"
	return desc
}

// Run implements cli.Command.
func (b *Username) Run(args []string) int {
	var events []map[string]interface{}
	err := validateInputArgs(args)
	if err != nil {
		b.UI.Error(fmt.Sprint(err))
		return 1
	} else {
		url := fmt.Sprintf("https://api.github.com/users/%v/events", args[0])
		resp, err := http.Get(url)
		if err != nil {
			b.UI.Error("something went wrong, please check the username")
			return 1
		}
		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	b.UI.Error("something went wrong, please try after sometime")
		// 	return 1
		// }
		// body := resp.Body.Close()
		b.UI.Output(resp.Status)
		json.NewDecoder(resp.Body).Decode(&events)
		for _, event := range events {
			fmt.Printf("%s: %s by %s\n",
				event["type"],
				event["repo"].(map[string]interface{})["name"],
				event["actor"].(map[string]interface{})["login"])
		}
		return 0
	}
}

// Synopsis implements cli.Command.
func (b *Username) Synopsis() string {
	synop := "This will give the use activity details. Usage: github-activity username <username>"
	return synop
}
