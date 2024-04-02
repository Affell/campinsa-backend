package planning

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"oui/models/user"
	"strings"

	"github.com/kataras/golog"
	"github.com/takuoki/clmconv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		golog.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		golog.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		golog.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func RetrievePlanning() (globalPlanning Planning) {

	spreadsheetId := os.Getenv("PLANNING_SPREADSHEET")
	ctx := context.Background()

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		golog.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		golog.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		golog.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	cached := make(map[string]user.User, 0)
	globalPlanning = make(Planning, 0)

	for iDay, day := range Days {

		resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, day+"!A1:1").Context(ctx).Do()
		if err != nil {
			golog.Error(err)
			return
		}

		if len(resp.Values) != 1 {
			fmt.Println("No data found.")
		} else {
			slots := make([]SheetSlot, 0)
			var current SheetSlot

			for j, col := range resp.Values[0] {
				s := col.(string)
				if s != "" {
					if current.Content != "" {
						current.End = j - 2
						slots = append(slots, current)
						current = SheetSlot{}
					}
					current.Content = s
					current.Start = j
				}
			}
			if current.Content != "" && current.End == 0 {
				current.End = -1
				slots = append(slots, current)
			}

			for _, slot := range slots {
				resp, err = srv.Spreadsheets.Values.Get(spreadsheetId, day+"!"+clmconv.Itoa(slot.Start)+"2:"+clmconv.Itoa(slot.End)+"35").Context(ctx).Do()
				if err != nil {
					golog.Error(err)
					return
				}

				for j := 0; j < len(resp.Values[0]); j += 2 {
					if resp.Values[0][j] == "Voitures" {
						continue
					}
					timeSlot := TimeSlot{
						Time:  slot.Content,
						Title: resp.Values[0][j].(string),
					}
					for i := 1; i < len(resp.Values); i++ {
						if j < len(resp.Values[i])-1 {
							lastname := strings.TrimSpace(resp.Values[i][j].(string))
							if resp.Values[i][j+1] != "" {
								firstname := strings.TrimSpace(resp.Values[i][j+1].(string))
								var u user.User
								if temp, ok := cached[firstname+lastname]; ok {
									u = temp
								} else {
									u = user.GetUserByName(firstname, lastname)
									if u.ID != 0 {
										cached[firstname+lastname] = u
									}
								}
								if u.ID != 0 {
									if _, ok := globalPlanning[u.ID]; !ok {
										globalPlanning[u.ID] = make(UserWeekPlanning, 0)
									}
									if _, ok := globalPlanning[u.ID][iDay]; !ok {
										globalPlanning[u.ID][iDay] = make(UserDayPlanning, 0)
									}

									temp := timeSlot
									if timeSlot.Title == "TAXI" {
										owner := strings.TrimSpace(resp.Values[i][j+3].(string)) + " " + strings.TrimSpace(resp.Values[i][j+2].(string))
										if owner != firstname+" "+lastname {
											temp.Additional = "Voiture : " + owner
										}
									}

									globalPlanning[u.ID][iDay] = append(globalPlanning[u.ID][iDay], temp)
								} else {
									golog.Errorf("Sheets: Unable to recognize %s %s", firstname, lastname)
								}
							} else {
								timeSlot = TimeSlot{
									Time:       timeSlot.Time,
									Title:      timeSlot.Title,
									Additional: lastname,
								}
							}

						}
					}
				}

			}
		}

	}

	f, err := os.OpenFile("planning.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		golog.Errorf("Unable to write planning: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(globalPlanning)

	return globalPlanning
}
