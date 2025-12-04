package cmds

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/cli"
)

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

		json.NewDecoder(resp.Body).Decode(&events)

		eventCounts := make(map[string]int)
		eventMessages := map[string]func(map[string]interface{}, string) (string, bool){
			"PushEvent": func(p map[string]interface{}, repoName string) (string, bool) {
				sizeRaw, _ := p["size"].(float64)
				size := int(sizeRaw)
				if size == 0 {
					size = 1
				}
				return fmt.Sprintf("Pushed %d commits to %s", size, repoName), true
			},
			"IssuesEvent": func(p map[string]interface{}, repoName string) (string, bool) {
				action, _ := p["action"].(string)
				if action == "opened" {
					return fmt.Sprintf("Opened a new issue in %s", repoName), true
				} else if action == "closed" {
					return fmt.Sprintf("Closed an issue in %s", repoName), true
				}
				return "", false
			},
			"WatchEvent": func(p map[string]interface{}, repoName string) (string, bool) {
				return fmt.Sprintf("Starred %s", repoName), true
			},
			"ForkEvent": func(p map[string]interface{}, repoName string) (string, bool) {
				forkee, _ := p["forkee"].(map[string]interface{})
				if forkeeName, ok := forkee["full_name"].(string); ok {
					return fmt.Sprintf("Forked %s", forkeeName), true
				}
				return fmt.Sprintf("Forked %s", repoName), true
			},
			"PullRequestEvent": func(p map[string]interface{}, repoName string) (string, bool) {
				action, _ := p["action"].(string)
				if action == "opened" {
					return fmt.Sprintf("Opened a pull request in %s", repoName), true
				} else if action == "closed" {
					return fmt.Sprintf("Closed a pull request in %s", repoName), true
				} else if action == "merged" {
					return fmt.Sprintf("Merged a pull request in %s", repoName), true
				}
				return "", false
			},
			"IssueCommentEvent": func(p map[string]interface{}, repoName string) (string, bool) {
				return fmt.Sprintf("Commented on an issue in %s", repoName), true
			},
		}

		// Process and count events
		for _, event := range events {
			eventType, ok := event["type"].(string)
			if !ok {
				continue
			}

			repo, ok := event["repo"].(map[string]interface{})
			if !ok {
				continue
			}
			repoName, ok := repo["name"].(string)
			if !ok {
				continue
			}

			payload, ok := event["payload"].(map[string]interface{})
			if !ok {
				continue
			}

			if handler, exists := eventMessages[eventType]; exists {
				if msg, shouldCount := handler(payload, repoName); shouldCount {
					// Use full message as key to group identical events
					key := msg
					eventCounts[key]++
				}
			}

		}
		// Sort by count (descending) then alphabetically
		type eventCount struct {
			message string
			count   int
		}
		var sortedEvents []eventCount
		for msg, count := range eventCounts {
			sortedEvents = append(sortedEvents, eventCount{msg, count})
		}
		sort.Slice(sortedEvents, func(i, j int) bool {
			if sortedEvents[i].count != sortedEvents[j].count {
				return sortedEvents[i].count > sortedEvents[j].count
			}
			return sortedEvents[i].message < sortedEvents[j].message
		})

		// Output
		fmt.Println("Output:")
		for _, ec := range sortedEvents {
			if ec.count > 1 {
				fmt.Printf("- %s (%d times)\n", ec.message, ec.count)
			} else {
				fmt.Printf("- %s\n", ec.message)
			}
		}

		return 0
	}
}

// Synopsis implements cli.Command.
func (b *Username) Synopsis() string {
	synop := "This will give the use activity details. Usage: github-activity username <username>"
	return synop
}
