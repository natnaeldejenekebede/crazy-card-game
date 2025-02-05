
// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math/rand"
// 	"net/http"
// 	"sync"
// 	"time"
// 	"github.com/gorilla/websocket"
// 	// "main/config"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// type Mes struct {
// 	Value string `json:"value"`
// 	Which string `json:"which"`
// 	Num   int    `json:"num"`
// }
// type Message struct {
// 	Value   string `json:"value"`
// 	Initial []Card `json:"initial"` // Nested struct
// 	Draw    []Card `json:"draw"`    // Nested struct
// }
// type remove struct {
// 	Value string `json:"value"`
// 	Cards []Card `json:"cards"`
// }

// type MoveMessage struct {
// 	Card Card `json:"card"` // Struct for the card played
// 	Draw bool `json:"draw"` // Whether the player wants to draw a card
// }

// type Player struct {
// 	ID        string
// 	Conn      *websocket.Conn
// 	Hand      []Card // Cards in the player's hand
// 	Drawlimit bool
// }

// type Card struct {
// 	Suit  string `json:"Suit"`
// 	Value string `json:"Value"`
// }

// type Game struct {
// 	Players   map[string]*Player
// 	Turn      string
// 	GameState string
// 	Deck      []Card
// 	PlayStack []Card
// 	mu        sync.Mutex // To handle concurrency
// }

// var game = Game{
// 	Players:   make(map[string]*Player),
// 	GameState: "waiting", // Game starts in 'waiting' state
// }

// var (
// 	state2    = 0
// 	stateA    = 0
// 	state2A   = 0
// 	state7    = ""
// 	first     = true
// 	drawlimit = false
// )

// func initializeDeck() {
// 	suits := []string{"Hearts", "Diamonds", "Clubs", "Spades"}
// 	values := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

// 	// Initialize deck with all cards
// 	for _, suit := range suits {
// 		for _, value := range values {
// 			game.Deck = append(game.Deck, Card{Suit: suit, Value: value})
// 		}
// 	}
// 	rand.NewSource(time.Now().UnixNano())
// 	// Shuffle the deck

// 	// rand.Seed(time.Now().UnixNano())
// 	for i := len(game.Deck) - 1; i > 0; i-- {
// 		j := rand.Intn(i + 1)
// 		game.Deck[i], game.Deck[j] = game.Deck[j], game.Deck[i]
// 	}
// }

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Error upgrading connection:", err)
// 		return
// 	}
// 	defer conn.Close()
// 	game.mu.Lock()
// 	playerID := fmt.Sprintf("Player-%d", len(game.Players)+1)
// 	fmt.Println(playerID)
// 	player := &Player{ID: playerID, Conn: conn}
// 	game.Players[playerID] = player
// 	game.mu.Unlock()

// 	fmt.Println(playerID, "connected")

// 	// Deal initial cards to the player (e.g., 5 cards)
// 	for i := 0; i < 5; i++ {
// 		card := drawCard()
// 		player.Hand = append(player.Hand, card)
// 	}
// 	if game.Players[playerID].ID == "Player-1" && len(game.Players[playerID].Hand) == 5 {
// 		card := drawCard()
// 		game.Players[playerID].Hand = append(game.Players[playerID].Hand, card)
// 	}
// 	message := Message{
// 		Value:   "initial",
// 		Initial: player.Hand,
// 	}

// 	jsonMsg, err := json.Marshal(message)
// 	// Send the JSON message to the WebSocket client
// 	err = conn.WriteMessage(websocket.TextMessage, jsonMsg)
// 	// conn.WriteJSON(message)
// 	// fmt.Println(jsonMsg)

// 	// Send initial game state to the player
// 	// conn.WriteJSON(game.GameState)
// 	if len(game.Players) >= 2 && game.GameState == "waiting" {
// 		startGame() // Start the game and set the game state to "in-progress"
// 	}

// 	for {
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Error reading message:", err)
// 			delete(game.Players, playerID)
// 			break
// 		}
// 		fmt.Printf("Received from %s: %s\n", playerID, string(msg))
// 		fmt.Println(game.GameState, game.Turn)
// 		// Handle game logic here (turn handling, card play validation, etc.)
// 		if game.GameState == "in-progress" && game.Turn == playerID && len(game.Deck) != 0 {
// 			handlePlayerMove(playerID, msg)
// 			me := game.Players[playerID]
// 			for _, j := range game.Players {
// 				if j.ID == playerID {
// 					continue
// 				}
// 				j.Conn.WriteJSON(Mes{Value: "oppounts", Which: me.ID, Num: len(me.Hand)})
// 				//conn.WriteJSON(Mes{Value: "oppounts",Num:len(j.Hand)})
// 			}
// 			if len(game.Players[playerID].Hand) == 0 {
// 				for _, j := range game.Players {
// 					if j.ID == playerID {
// 						j.Conn.WriteJSON(Mes{Value: "won", Num: len(me.Hand)})
// 						game.GameState = "waiting"
// 						continue
// 					}
// 					j.Conn.WriteJSON(Mes{Value: "loss", Num: len(me.Hand)})
// 					game.GameState = "waiting"
// 					//conn.WriteJSON(Mes{Value: "oppounts",Num:len(j.Hand)})
// 				}

// 			}

// 		}
// 	}
// }

// // Function to start the game
// func startGame() {
// 	if len(game.Players) >= 2 { // Ensure there are enough players
// 		game.GameState = "in-progress" // Set the game state to "in-progress"
// 		game.Turn = "Player-1"         // Decide the starting player
// 		fmt.Println("Game started!")
// 	} else {
// 		fmt.Println("Not enough players to start the game")
// 	}
// }

// func drawCard() Card {
// 	// Lock the mutex to ensure only one player can draw a card at a time
// 	// game.mu.Lock()
// 	// defer game.mu.Unlock()

// 	// Check if the deck has any cards left
// 	if len(game.Deck) == 0 {
// 		fmt.Println("No cards left in the deck")
// 		return Card{} // Return an empty card if no cards are left
// 	}

// 	// Draw a card from the deck
// 	card := game.Deck[0]
// 	game.Deck = game.Deck[1:] // Remove the card from the deck
// 	if len(game.Deck) == 1 {
// 		var rem remove
// 		rem.Value = "empty"

// 		game.Players[game.Turn].Conn.WriteJSON(rem)

// 	}
// 	// Return the drawn card
// 	return card
// }

// func handlePlayerMove(playerID string, msg []byte) {
// 	fmt.Println("entered to hanldle player")
// 	game.mu.Lock()
// 	defer game.mu.Unlock()
// 	var message MoveMessage
// 	// Decode the incoming message into a Card struct

// 	err := json.Unmarshal(msg, &message)
// 	var cardPlayed Card = message.Card
// 	fmt.Println(cardPlayed, message)
// 	var draw bool = message.Draw
// 	if err != nil {
// 		fmt.Println("Error decoding card:", err)
// 		return
// 	}
// 	fmt.Println("entered to check player")
// 	// Check if it's the player's turn
// 	// if playerID != game.Turn {
// 	// 	// Send error message back to player if it's not their turn
// 	// 	for _, player := range game.Players {
// 	// 		player.Conn.WriteJSON(fmt.Sprintf("%s, it's not your turn!", playerID))
// 	// 	}
// 	// 	return
// 	// }
// 	//fmt.Println(game.Players)
// 	// player, exists := game.Players["Player-1"]
// 	// if !exists {
// 	// 	fmt.Println("Player-1 does not exist in game.Players")
// 	// } else {

// 	// 	fmt.Println("Player-1's hand size:", len(player.Hand))

// 	// }
// 	fmt.Println(playerID)
// 	if playerID == "Player-1" && first {
// 		if draw {
// 			return
// 		}
// 		fmt.Println("entered to first player")
// 		if newvailed(cardPlayed, playerID) {
// 			game.PlayStack = append(game.PlayStack, cardPlayed)
// 			removeCardFromHand(playerID, cardPlayed)
// 			first = false
// 			var rem remove
// 			rem.Value = "remove"
// 			rem.Cards = append(rem.Cards, cardPlayed)
// 			game.Players[playerID].Conn.WriteJSON(rem)
// 			fmt.Println("sent to first player")
// 			fmt.Println(cardPlayed, playerID)
// 			nextTurn(game.Turn)
// 			tellothers(cardPlayed, playerID)
// 			fmt.Println("entered to finsh ")
// 		}

// 	} else {
// 		// Validate and apply the card play according to game rules
// 		if draw {
// 			for _, j := range game.Players {
// 				if j.ID == game.Turn {
// 					continue
// 				}
// 				j.Drawlimit = false
// 			}
// 			currentCard := game.PlayStack[len(game.PlayStack)-1]
// 			if state7 == "add" { //currentCard.Value == "7" &&
// 				if len(game.Players) > 2 {
// 					fmt.Println("entered to check 7", state7)
// 					return

// 				} else {
// 					fmt.Println("entered to ch", state7)
// 					return
// 				}

// 			}
// 			if currentCard.Value == "5" {

// 			}
// 			if game.Players[game.Turn].Drawlimit {
// 				fmt.Println("enered draw limit", game.Players[game.Turn])
// 				game.Players[game.Turn].Drawlimit = false
// 				nextTurn(playerID)
// 				return
// 			}
// 			fmt.Println("draw reached")
// 			if state2 >= 1 {
// 				var rem remove
// 				rem.Value = "add"

// 				fmt.Println(state2)
// 				card := drawCard()
// 				rem.Cards = append(rem.Cards, card)
// 				game.Players[game.Turn].Hand = append(game.Players[game.Turn].Hand, card)

// 				game.Players[playerID].Conn.WriteJSON(rem)
// 				fmt.Println(game.Players[game.Turn].Hand)
// 				// nextTurn(playerID)

// 				state2 -= 1
// 				if state2 == 0 {
// 					fmt.Println("state becomes zero")
// 					game.Players[playerID].Drawlimit = true
// 				}

// 			} else {
// 				fmt.Println("to lecture")
// 				card := drawCard()
// 				game.Players[game.Turn].Hand = append(game.Players[game.Turn].Hand, card)
// 				var rem remove
// 				rem.Value = "add"
// 				rem.Cards = append(rem.Cards, card)
// 				game.Players[playerID].Conn.WriteJSON(rem)
// 				game.Players[playerID].Drawlimit = true
// 				fmt.Println(game.Players[playerID])

// 				// nextTurn(playerID)
// 			}

// 		} else {
// 			if validMove(cardPlayed, playerID) {
// 				fmt.Println("entered to vailedmoves")
// 				game.Players[game.Turn].Drawlimit = false
// 				game.Players[playerID].Drawlimit = false

// 				// Add the played card to the play stack
// 				// if cardPlayed.Value == "5" && len(game.Players) > 2 {
// 				// 	game.Players[game.Turn].Drawlimit = true
// 				// }
// 				game.PlayStack = append(game.PlayStack, cardPlayed)

// 				// Remove the card from the player's hand
// 				removeCardFromHand(playerID, cardPlayed)
// 				var rem remove
// 				rem.Value = "remove"
// 				rem.Cards = append(rem.Cards, cardPlayed)
// 				game.Players[playerID].Conn.WriteJSON(rem)
// 				// var rem remove
// 				// rem.value = "remove"
// 				// rem.cards = []Card{cardPlayed} // Make sure cards is properly initialized
// 				// game.Players[playerID].Conn.WriteJSON(rem)
// 				tellothers(cardPlayed, playerID)
// 				// Move to the next turn
// 				nextTurn(game.Turn)
// 			}
// 		}

// 	}

// }
// func newvailed(card Card, playerID string) bool {
// 	if card.Value == "8" || card.Value == "J" {
// 		// changeSuit(card)
// 		//how to change the suit of the play
// 		return true
// 	}
// 	if card.Value == "2" {
// 		//let the next player pick 2 extra cards
// 		state2 = 2
// 		game.GameState = "in-progress"
// 		return true
// 	}
// 	if card.Value == "5" {
// 		//jumpe the next player
// 		nextTurn(playerID)
// 		return true
// 	}
// 	if card.Value == "7" {
// 		//handle to play all simlar suits at once
// 		handle7Card(card, playerID)
// 		state7 = "add"
// 		first = false
// 		// Handle reverse direction if only , depending on the game state (not implemented in this example)
// 		return false
// 	}
// 	if card.Value == "A" && card.Suit == "Spades" {
// 		//punsih the next player to get 5 cards from teh deck
// 		draw5CardsFromDeck()
// 		return true

// 	}
// 	return true
// }

// // Example validation function
// func validMove(card Card, playerID string) bool {

// 	// Check if the card matches the current suit or number
// 	currentCard := game.PlayStack[len(game.PlayStack)-1] // Last played card
// 	if currentCard.Value == "2" && card.Value == "2" {
// 		state2 += 2
// 		// for i:=0;i<state2;i++{
// 		// 	game.mu.Lock()
// 		// 	card:=drawCard()
// 		// 	game.Players[game.Turn].Hand = append(game.Players[game.Turn].Hand, card)
// 		// 	game.mu.Unlock()
// 		// }
// 		return true
// 	}
// 	if currentCard.Value == "2" && card.Value != "2" {
// 		if state2 != 0 {
// 			return false
// 		}
// 		// return true
// 	}
// 	if state7 == "add" && card.Suit == currentCard.Suit {
// 		fmt.Println("entered to add becouse of 7")
// 		handle7Card(card, playerID)
// 		return false

// 	}
// 	if state7 == "add" && card.Suit != currentCard.Suit {
// 		if currentCard.Value == "7" {
// 			fmt.Println("entered to reverse becouse of 7")
// 			reverseGamePlayers(playerID)
// 			game.Players[game.Turn].Drawlimit = false
// 			state7 = ""
// 			return false
// 		}
// 		fmt.Println("entered finshied addingof 7")
// 		state7 = ""
// 		nextTurn(playerID)
// 		return false
// 	}

// 	// If the card is a wild card (8 or J), it can be played on any card
// 	if card.Value == "8" || card.Value == "J" {
// 		return changeSuit(card, currentCard)

// 		//how to change the suit of the play

// 	}

// 	// Check if the card matches the suit or value of the current card
// 	if card.Suit == currentCard.Suit || card.Value == currentCard.Value {
// 		if card.Value == "2" {
// 			//let the next player pick 2 extra cards
// 			state2 += 2
// 			game.GameState = "in-progress"
// 			return true
// 		}
// 		if card.Value == "5" {
// 			//jumpe the next player
// 			nextTurn(playerID)
// 			return true
// 		}
// 		if card.Value == "7" {
// 			fmt.Println("entered proper 7")
// 			//handle to play all simlar suits at once
// 			handle7Card(card, playerID)
// 			state7 = "add"
// 			// Handle reverse direction if only , depending on the game state (not implemented in this example)
// 			return false
// 		}
// 		if card.Value == "A" && card.Suit == "Spades" {
// 			//punsih the next player to get 5 cards from teh deck
// 			draw5CardsFromDeck()
// 			return true

// 		}
// 		//noramll
// 		return true
// 	}

// 	return false
// }
// func tellothers(cardPlayed Card, playerID string) {
// 	fmt.Println("entered to tell others", game.Players)
// 	// game.mu.Lock()
// 	// defer game.mu.Unlock()

// 	fmt.Println("entered first", game.Players)
// 	for i, j := range game.Players {
// 		if i == playerID {
// 			continue
// 		}
// 		fmt.Println("entered")
// 		var mes Message
// 		mes.Value = "top"
// 		mes.Initial = append(mes.Initial, cardPlayed)
// 		fmt.Println("entered", mes.Initial)
// 		j.Conn.WriteJSON(mes)
// 		fmt.Println("finsished sending")
// 	}

// }
// func changeSuit(card Card, currentcard Card) bool {
// 	if currentcard.Value == "J" || currentcard.Value == "8" {
// 		if currentcard.Value == card.Value || currentcard.Suit == card.Suit {
// 			return true
// 		}
// 		game.PlayStack = append(game.PlayStack, currentcard)

// 		// Remove the card from the player's hand

// 		removeCardFromHand(game.Turn, card)
// 		fmt.Println(game.Players[game.Turn].Hand)
// 		var rem remove
// 		rem.Value = "change"
// 		rem.Cards = append(rem.Cards, card)
// 		game.Players[game.Turn].Conn.WriteJSON(rem)
// 		// fmt.Println(game.Players[game.Turn].Hand)
// 		// var add remove
// 		// add.Value = "add"
// 		// add.Cards = append(add.Cards, currentcard)
// 		// game.Players[game.Turn].Conn.WriteJSON(add)
// 		// fmt.Println(game.Players[game.Turn].Hand)
// 		// var rem2 remove
// 		// rem2.Value = "remove"
// 		// rem2.Cards = append(rem2.Cards, currentcard)
// 		// game.Players[game.Turn].Conn.WriteJSON(rem2)
// 		// fmt.Println(game.Players[game.Turn].Hand)
// 		// var rem remove
// 		// rem.value = "remove"
// 		// rem.cards = []Card{cardPlayed} // Make sure cards is properly initialized
// 		// game.Players[playerID].Conn.WriteJSON(rem)
// 		tellothers(currentcard, game.Turn)
// 		// Move to the next turn
// 		nextTurn(game.Turn)
// 		return false

// 	}
// 	// Update the game state to reflect the new suit in play
// 	game.GameState = "in-progress" //fmt.Sprintf("new-suit-%s", card.Suit)
// 	fmt.Printf("Suit changed to %s\n", card.Suit)
// 	return true
// }

// // func skipNextPlayer(PlayerID string) {
// // 	fmt.Println(game.Turn)
// // 	nextTurn(PlayerID)
// // 	// nextTurn(game.Turn)
// // 	fmt.Println(game.Turn)
// // 	// Implement logic to skip the next player

// // 	fmt.Println("Skipping the next player.")
// // 	// You would need to update the game turn or player list to skip the next player
// // }

// func handle7Card(card Card, playerID string) {
// 	fmt.Println("entered 7")
// 	// game.mu.Lock()
// 	game.PlayStack = append(game.PlayStack, card)
// 	removeCardFromHand(game.Turn, card)
// 	var rem remove
// 	rem.Value = "remove"
// 	rem.Cards = []Card{card} // Make sure cards is properly initialized
// 	game.Players[playerID].Conn.WriteJSON(rem)
// 	tellothers(card, playerID)
// 	game.Turn = game.Turn
// 	// game.mu.Unlock()

// 	// Remove the card from the player's hand

// 	// Move to the next turn

// }

// // Function to handle drawing 5 cards from the deck when Ace of Spades is played
// func draw5CardsFromDeck() {
// 	fmt.Println("Next player must draw 5 cards from the deck.")
// 	// Modify the game state and the player’s hand to simulate drawing cards
// }

// // Remove the played card from the player's hand
// func removeCardFromHand(playerID string, card Card) {
// 	player := game.Players[playerID]
// 	for i, c := range player.Hand {
// 		if c == card {
// 			// Remove the card from the player's hand
// 			player.Hand = append(player.Hand[:i], player.Hand[i+1:]...)
// 			break
// 		}
// 	}
// }

// func nextTurn(currentPlayerID string) {
// 	fmt.Println("enterd tot urn turn")
// 	// Rotate to the next player for the turn
// 	playerIDs := []string{}
// 	for playerID := range game.Players {
// 		playerIDs = append(playerIDs, playerID)
// 	}

// 	// Find the index of the current player
// 	for i, id := range playerIDs {
// 		if id == currentPlayerID {
// 			// Set the next player as the current turn
// 			nextPlayerID := playerIDs[(i+1)%len(playerIDs)]
// 			// game.mu.Lock()
// 			game.Turn = nextPlayerID
// 			// game.mu.Unlock()
// 			fmt.Println("turn", game.Turn)
// 			break
// 		}
// 	}
// 	// if {
// 	// 	game.mu.Lock()

// 	// 	card:=drawCard()
// 	// 	game.Players[game.Turn].Hand = append(game.Players[game.Turn].Hand, card)
// 	// 	game.mu.Unlock()
// 	// }

// 	// Notify all players of the next turn
// 	for _, player := range game.Players {
// 		player.Conn.WriteJSON(fmt.Sprintf("It's %s's turn", game.Turn))
// 	}
// }

// func reverseGamePlayers(playerID string) {
// 	// Step 1: Create a slice of playerIDs
// 	playerIDs := []string{}
// 	for playerID := range game.Players {
// 		playerIDs = append(playerIDs, playerID)
// 	}

// 	// Step 2: Reverse the playerIDs slice
// 	for i, j := 0, len(playerIDs)-1; i < j; i, j = i+1, j-1 {
// 		playerIDs[i], playerIDs[j] = playerIDs[j], playerIDs[i]
// 	}

// 	// Step 3: Create a new map for reversed players
// 	reversedPlayers := make(map[string]*Player)
// 	for _, playerID := range playerIDs {
// 		reversedPlayers[playerID] = game.Players[playerID]
// 	}

// 	// Step 4: Assign the reversed map to game.Players
// 	game.Players = reversedPlayers
// 	nextTurn(playerID)
// 	// Print out the reversed player order (for testing)
// 	fmt.Println("Reversed game.Players:", game.Players)
// }

// func main() {
// 	initializeDeck() // Initialize and shuffle the deck
// 	http.HandleFunc("/ws", handleConnections)
// 	fmt.Println("Server started on :8080")
// 	http.ListenAndServe(":8080", nil)
// }
