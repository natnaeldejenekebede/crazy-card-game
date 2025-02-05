// import React, { useState, useEffect, useRef } from "react";
// import Card from "./Card";
// import OpponentCards from "./OpponentCards";
// import "./CrazyCardGame.css";

// const wsUrl = "ws://192.168.100.5:8080/ws";

// const CrazyCardGame = () => {
//   const hello = useRef([]);
//   const [hand, setHand] = useState([]);
//   const [topCard, setTopCard] = useState(null);
//   const [gameState, setGameState] = useState("in-progress");
//   const [turn, setTurn] = useState("");
//   const [ws, setWs] = useState(null);
//   const [opponents, setOpponents] = useState([]);
//   const [opponents2, setOpponents2] = useState([]);

//   useEffect(() => {
//     const socket = new WebSocket(wsUrl);

//     socket.onopen = () => {
//       console.log("Connected to the WebSocket server.");
//     };

//     socket.onmessage = (event) => {
//       const message = JSON.parse(event.data);
//       console.log(message);

//       if (message.value === "initial") {
//         setHand(message.initial);
//       } else if (message.value === "remove") {
//         const firstCard = message.cards[0];
//         const updatedHand = hello.current.filter(card => !(card.Suit === firstCard.Suit && card.Value === firstCard.Value));
//         setHand(updatedHand);
//         setTopCard(message.cards[0]);
//       } else if (message.value === "add") {
//         setHand((prevHand) => [...prevHand, ...message.cards]);
//       } else if (message.value === "top") {
//         setTopCard(message.initial[0]);
//       } else if (message.value === "change") {
//         const firstCard = message.cards[0];
//         const updatedHand = hello.current.filter(card => !(card.Suit === firstCard.Suit && card.Value === firstCard.Value));
//         setHand(updatedHand);
//       }else if (message.value==="empty"){
//         alert("empty deck")
//       } 
//       else if (message.value === "oppounts") {
//         if(!opponents.which){
//           setOpponents(message)
//         }
//         else if (!opponents2.which){
//           setOpponents2(message)
//         }
//         if (opponents.which || opponents2.which){
//           if (message.which==opponents.which){
//             setOpponents(message)
//           }
//           if(message.which==opponents2.which){
//             setOpponents2(message)
//           }

//         }
          
       
        
//       } else if (message.value ==="won"){
//         alert("You Won")
//       }else if(message.value=== "loss"){
//         alert("You lost")
//       }

//       if (message.value === "game-state") {
//         setGameState(message.gameState);
//         setTurn(message.turn);
//       }
//     };

//     socket.onclose = () => {
//       console.log("Disconnected from the WebSocket server.");
//     };

//     setWs(socket);

//     return () => {
//       socket.close();
//     };
//   }, []);

//   const playCard = (index) => {
//     if (ws && gameState === "in-progress") {
//       hello.current = hand;
//       const cardToPlay = hand[index];
//       const moveMessage = { card: cardToPlay, draw: false };
//       ws.send(JSON.stringify(moveMessage));
//     }
//   };

//   const drawCard = () => {
//     if (ws && gameState === "in-progress") {
//       const drawMessage = { card: null, draw: true };
//       ws.send(JSON.stringify(drawMessage));
//     }
//   };

//   return (
//     <div className="game-container">
//       <h1 className="game-title">Crazy Card Game</h1>
//       {/* {if (opponents.message)} */}
//       <div style={{ display: 'flex', gap: '10px' }}>
//       <OpponentCards opponents={opponents} />
//       <OpponentCards opponents={opponents2} />
//       </div>
     
     

//       <div className="top-card-container">
//         <h2>Top Card</h2>
//         {topCard && <Card suit={topCard.Suit} value={topCard.Value} onPlay={() => {}} />}
//       </div>

//       <div className="hand-container">
//         {hand.map((card, index) => (
//           <Card key={index} suit={card.Suit} value={card.Value} onPlay={() => playCard(index)} />
//         ))}
//       </div>

//       <button className="draw-button" onClick={drawCard} disabled={hand.length >= 10}>
//         Draw Card ({hand.length}/10)
//       </button>
//     </div>
//   );
// };

// export default CrazyCardGame;





import React, { useState, useEffect, useRef } from "react";
import Card from "./Card";
import OpponentCards from "./OpponentCards";
import "./CrazyCardGame.css";
import { useParams } from "react-router-dom";

const CrazyCardGame = () => {
  const hello = useRef([]);
  const isReconnecting = useRef(false);  // Prevents multiple reconnection attempts
  const [hand, setHand] = useState([]);
  const [topCard, setTopCard] = useState(null);
  const [gameState, setGameState] = useState("in-progress");
  const [turn, setTurn] = useState("");
  const [ws, setWs] = useState(null);
  const [opponents, setOpponents] = useState([]);
  const [opponents2, setOpponents2] = useState([]);
  const [connected, setConnected] = useState(false);  // Ensures rendering happens only after first message
  const { suit } = useParams();

  useEffect(() => {
    const fetchAvailableIP = async () => {
      try {
        console.log("Fetching new IP...");
        const response = await fetch("http://192.168.100.5:8083/reconnect");
        if (!response.ok) {
          throw new Error("Failed to fetch new IP");
        }
        const data = await response.json();
        return data.newip;
      } catch (error) {
        console.error("Error fetching IP:", error);
        return null;
      }
    };

    const wsConnect = (wsUrl) => {
      const socket = new WebSocket(wsUrl);
      console.log("Connecting to WebSocket:", wsUrl);

      socket.onopen = () => {
        console.log("Connected to WebSocket server.");
        setConnected(false); // Reset so we wait for the first message before rendering
      };

      socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log("Received message:", message);

        if (!connected) {
          setConnected(true); // Now, we allow the game to render
        }

        switch (message.value) {
          case "initial":
            setHand(message.initial);
            break;
          case "remove":
            const firstCard = message.cards[0];
            setHand((prevHand) => prevHand.filter(card => !(card.Suit === firstCard.Suit && card.Value === firstCard.Value)));
            setTopCard(message.cards[0]);
            break;
          case "add":
            setHand((prevHand) => [...prevHand, ...message.cards]);
            break;
          case "top":
            setTopCard(message.initial[0]);
            break;
          case "change":
            const changedCard = message.cards[0];
            setHand((prevHand) => prevHand.filter(card => !(card.Suit === changedCard.Suit && card.Value === changedCard.Value)));
            break;
          case "empty":
            alert("Deck is empty");
            break;
          case "opponents":
            if (!opponents.which) setOpponents(message);
            else if (!opponents2.which) setOpponents2(message);
            else {
              if (message.which === opponents.which) setOpponents(message);
              if (message.which === opponents2.which) setOpponents2(message);
            }
            break;
          case "won":
            alert("You Won!");
            break;
          case "loss":
            alert("You Lost!");
            break;
          case "game-state":
            setGameState(message.gameState);
            setTurn(message.turn);
            break;
          default:
            console.warn("Unknown message type:", message.value);
        }
      };

      socket.onclose = async () => {
        console.log("WebSocket disconnected.");

        if (!isReconnecting.current) {
          isReconnecting.current = true; // Prevent multiple reconnections

          const newIp = await fetchAvailableIP();
          if (newIp) {
            const newWsUrl = `ws://${newIp}/restart`;
            console.log("Reconnecting to new IP:", newWsUrl);
            wsConnect(newWsUrl);

            // Send the username after reconnecting
            socket.onopen = () => {
              const moveMessage = { "username":sessionStorage.getItem("username")};
              console.log("entered to send     tt")
//             socket.send(JSON.stringify(moveMessage));
              // const username = sessionStorage.getItem("username");
              socket.send(JSON.stringify(moveMessage));
              isReconnecting.current = false; // Reset after successful reconnection
            };
          }
        }
      };

      socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      setWs(socket);
    };

    const initialIp = localStorage.getItem("ip");
    if (!initialIp) {
      console.error("WebSocket URL not found in localStorage.");
      return;
    }

    wsConnect(`ws://${initialIp}/ws`);

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, []);

  const playCard = (index) => {
    if (ws && gameState === "in-progress") {
      hello.current = hand;
      const cardToPlay = hand[index];
      const username = sessionStorage.getItem("username");
      const moveMessage = { card: cardToPlay, draw: false, username };
      ws.send(JSON.stringify(moveMessage));
    }
  };

  const drawCard = () => {
    if (ws && gameState === "in-progress") {
      const username = sessionStorage.getItem("username");
      const drawMessage = { card: null, draw: true, username };
      ws.send(JSON.stringify(drawMessage));
    }
  };

  return (
    <div className="game-container">
      <h1>Welcome to the {suit.charAt(0).toUpperCase() + suit.slice(1)} Room!</h1>
      <h1>Welcome, {sessionStorage.getItem("username")}</h1>
      <h1 className="game-title">Crazy Card Game</h1>

      {!connected ? (
        <p>Waiting for game data...</p>
      ) : (
        <>
          <div style={{ display: 'flex', gap: '10px' }}>
            <OpponentCards opponents={opponents} />
            <OpponentCards opponents={opponents2} />
          </div>

          <div className="top-card-container">
            <h2>Top Card</h2>
            {topCard && <Card suit={topCard.Suit} value={topCard.Value} onPlay={() => {}} />}
          </div>

          <div className="hand-container">
            {hand.map((card, index) => (
              <Card key={index} suit={card.Suit} value={card.Value} onPlay={() => playCard(index)} />
            ))}
          </div>

          <button className="draw-button" onClick={drawCard} disabled={hand.length >= 10}>
            Draw Card ({hand.length}/10)
          </button>
        </>
      )}
    </div>
  );
};

export default CrazyCardGame;








































// import React, { useState, useEffect, useRef } from "react";
// import Card from "./Card";
// import OpponentCards from "./OpponentCards";
// import "./CrazyCardGame.css";
// import { useParams } from "react-router-dom";

// // const wsUrl = "ws://192.168.100.5:8080/ws";
// const wsUrl = localStorage.getItem("ip"); 

// const CrazyCardGame = () => {

//   const hello = useRef([]);
//   const hello2 = useRef([]);
//   const [hand, setHand] = useState([]);
//   const [topCard, setTopCard] = useState(null);
//   const [gameState, setGameState] = useState("in-progress");
//   const [turn, setTurn] = useState("");
//   const [ws, setWs] = useState(null);
//   const [opponents, setOpponents] = useState([]);
//   const [opponents2, setOpponents2] = useState([]);
//   const {suit}=useParams()
//   console.log(suit)

//   useEffect(() => {
//     const fetchAvailableIP = async () => {
//       try {
//         console.log("entered to refetch")
//         const response = await fetch("http://192.168.100.5:8083/reconnect");  // Request from authentication server
//         if (!response.ok) {
//           throw new Error("Failed to fetch available IP from auth server");
//         }
//         const data = await response.json();
//         console.log(data)
//         return data.newip;  // Assuming the response returns an object with an `ip` field
//       } catch (error) {
//         console.error("Error fetching IP:", error);
//         return null;
//       }
//     };
  
//     const wsConnect = (wsUrl) => {
//       const socket = new WebSocket(wsUrl);
//       console.log("Connecting to WebSocket:", wsUrl);
  
//       socket.onopen = () => {
//         console.log("Connected to the WebSocket server.");
//         console.log(sessionStorage.getItem("username"))
        
        

//       };
  
//       socket.onmessage = (event) => {
//         socket.onmessage = (event) => {
//                 const message = JSON.parse(event.data);
//                 console.log(message);
          
//                 if (message.value === "initial") {
//                   setHand(message.initial);
//                 } else if (message.value === "remove") {
//                   const firstCard = message.cards[0];
//                   const updatedHand = hello.current.filter(card => !(card.Suit === firstCard.Suit && card.Value === firstCard.Value));
//                   setHand(updatedHand);
//                   setTopCard(message.cards[0]);
//                 } else if (message.value === "add") {
//                   setHand((prevHand) => [...prevHand, ...message.cards]);
//                 } else if (message.value === "top") {
//                   setTopCard(message.initial[0]);
//                 } else if (message.value === "change") {
//                   const firstCard = message.cards[0];
//                   const updatedHand = hello.current.filter(card => !(card.Suit === firstCard.Suit && card.Value === firstCard.Value));
//                   setHand(updatedHand);
//                 }else if (message.value==="empty"){
//                   alert("empty deck")
//                 } 
//                 else if (message.value === "oppounts") {
//                   if(!opponents.which){
//                     setOpponents(message)
//                   }
//                   else if (!opponents2.which){
//                     setOpponents2(message)
//                   }
//                   if (opponents.which || opponents2.which){
//                     if (message.which==opponents.which){
//                       setOpponents(message)
//                     }
//                     if(message.which==opponents2.which){
//                       setOpponents2(message)
//                     }
          
//                   }
                    
                 
                  
//                 } else if (message.value ==="won"){
//                   alert("You Won")
//                 }else if(message.value=== "loss"){
//                   alert("You lost")
//                 }
          
//                 if (message.value === "game-state") {
//                   setGameState(message.gameState);
//                   setTurn(message.turn);
//                 }
//               };
          
//       };
//       hello2.current=true
//       socket.onclose = async () => {
        
//         console.log("Disconnected from the WebSocket server.");
//         // Fetch new IP address from auth server
//        if (hello2.current===true){
//         const newIp = await fetchAvailableIP();
//         if (newIp) {
//           const newWsUrl = `ws://${newIp}/restart`;
//           console.log("Reconnecting to new IP:", newWsUrl);
          
//           wsConnect(newWsUrl); 
//           if (hello2.current){
//             const moveMessage = { "username":sessionStorage.getItem("username")};
//             socket.send(JSON.stringify(moveMessage));
//           } // Reconnect with the new IP
//           hello2.current=false
//         }
//        }
       
       
        
//       };
  
//       socket.onerror = (error) => {
//         console.error("WebSocket error:", error);
//       };
  
//       setWs(socket);  // Save the socket to the state or a ref to handle it globally
//     };
  
//     const wsUrl = localStorage.getItem("ip");
//     if (!wsUrl) {
//       console.error("WebSocket URL not found in localStorage.");
//       return;
//     }
  
//     const initialWsUrl = `ws://${wsUrl}/ws`;
//     wsConnect(initialWsUrl);
  
//     return () => {
//       if (ws) {
//         ws.close();  // Cleanup on component unmount
//       }
//     };
//   }, []);
  
//   const playCard = (index) => {
//     if (ws && gameState === "in-progress") {
//       hello.current = hand;
//       const cardToPlay = hand[index];
//       const username = sessionStorage.getItem("username");
//       const moveMessage = { card: cardToPlay, draw: false ,username};
//       ws.send(JSON.stringify(moveMessage));
//     }
//   };

//   const drawCard = () => {
//     if (ws && gameState === "in-progress") {
//       const username = sessionStorage.getItem("username");
//       const drawMessage = { card: null, draw: true ,username};
//       ws.send(JSON.stringify(drawMessage));
//     }
//   };

//   return (
//     <div className="game-container">
//       <h1>Welcome to the {suit.charAt(0).toUpperCase() + suit.slice(1)} Room!</h1>
//       <h1>wllcom {sessionStorage.getItem("username")}</h1>
//       <h1 className="game-title">Crazy Card Game</h1>
//       {/* {if (opponents.message)} */}
//       <div style={{ display: 'flex', gap: '10px' }}>
//       <OpponentCards opponents={opponents} />
//       <OpponentCards opponents={opponents2} />
//       </div>
     
     

//       <div className="top-card-container">
//         <h2>Top Card</h2>
//         {topCard && <Card suit={topCard.Suit} value={topCard.Value} onPlay={() => {}} />}
//       </div>

//       <div className="hand-container">
//         {hand.map((card, index) => (
//           <Card key={index} suit={card.Suit} value={card.Value} onPlay={() => playCard(index)} />
//         ))}
//       </div>

//       <button className="draw-button" onClick={drawCard} disabled={hand.length >= 10}>
//         Draw Card ({hand.length}/10)
//       </button>
//     </div>
//   );
// };

// export default CrazyCardGame;
























// import React, { useState, useEffect, useRef } from "react";
// import Card from "./Card";
// import OpponentCards from "./OpponentCards";
// import "./CrazyCardGame.css";
// import { useParams } from "react-router-dom";

// // const wsUrl = "ws://192.168.100.5:8080/ws";
// const wsUrl = localStorage.getItem("ip"); 

// const CrazyCardGame = () => {

//   const hello = useRef([]);
//   const [hand, setHand] = useState([]);
//   const [topCard, setTopCard] = useState(null);
//   const [gameState, setGameState] = useState("in-progress");
//   const [turn, setTurn] = useState("");
//   const [ws, setWs] = useState(null);
//   const [opponents, setOpponents] = useState([]);
//   const [opponents2, setOpponents2] = useState([]);
//   const {suit}=useParams()
//   console.log(suit)

//   useEffect(() => {
//     // let  wsUrl = localStorage.getItem("ip");
//     let wsUrl = localStorage.getItem("ip");
//     if (!wsUrl) {
//       console.error("WebSocket URL not found in localStorage.");
//       return;
//     }
   
//  wsUrl = `ws://${wsUrl}/ws`;
//     const socket = new WebSocket(wsUrl);
//     console.log(wsUrl)

//     socket.onopen = () => {
//       console.log("Connected to the WebSocket server.");
//     };

//     socket.onmessage = (event) => {
//       const message = JSON.parse(event.data);
//       console.log(message);

//       if (message.value === "initial") {
//         setHand(message.initial);
//       } else if (message.value === "remove") {
//         const firstCard = message.cards[0];
//         const updatedHand = hello.current.filter(card => !(card.Suit === firstCard.Suit && card.Value === firstCard.Value));
//         setHand(updatedHand);
//         setTopCard(message.cards[0]);
//       } else if (message.value === "add") {
//         setHand((prevHand) => [...prevHand, ...message.cards]);
//       } else if (message.value === "top") {
//         setTopCard(message.initial[0]);
//       } else if (message.value === "change") {
//         const firstCard = message.cards[0];
//         const updatedHand = hello.current.filter(card => !(card.Suit === firstCard.Suit && card.Value === firstCard.Value));
//         setHand(updatedHand);
//       }else if (message.value==="empty"){
//         alert("empty deck")
//       } 
//       else if (message.value === "oppounts") {
//         if(!opponents.which){
//           setOpponents(message)
//         }
//         else if (!opponents2.which){
//           setOpponents2(message)
//         }
//         if (opponents.which || opponents2.which){
//           if (message.which==opponents.which){
//             setOpponents(message)
//           }
//           if(message.which==opponents2.which){
//             setOpponents2(message)
//           }

//         }
          
       
        
//       } else if (message.value ==="won"){
//         alert("You Won")
//       }else if(message.value=== "loss"){
//         alert("You lost")
//       }

//       if (message.value === "game-state") {
//         setGameState(message.gameState);
//         setTurn(message.turn);
//       }
//     };

//     socket.onclose = () => {
//       console.log("Disconnected from the WebSocket server.");
//     };

//     setWs(socket);

//     return () => {
//       socket.close();
//     };
//   }, []);

//   const playCard = (index) => {
//     if (ws && gameState === "in-progress") {
//       hello.current = hand;
//       const cardToPlay = hand[index];
//       const username = sessionStorage.getItem("username");
//       const moveMessage = { card: cardToPlay, draw: false ,username};
//       ws.send(JSON.stringify(moveMessage));
//     }
//   };

//   const drawCard = () => {
//     if (ws && gameState === "in-progress") {
//       const username = sessionStorage.getItem("username");
//       const drawMessage = { card: null, draw: true ,username};
//       ws.send(JSON.stringify(drawMessage));
//     }
//   };

//   return (
//     <div className="game-container">
//       <h1>Welcome to the {suit.charAt(0).toUpperCase() + suit.slice(1)} Room!</h1>
//       <h1>wllcom {sessionStorage.getItem("username")}</h1>
//       <h1 className="game-title">Crazy Card Game</h1>
//       {/* {if (opponents.message)} */}
//       <div style={{ display: 'flex', gap: '10px' }}>
//       <OpponentCards opponents={opponents} />
//       <OpponentCards opponents={opponents2} />
//       </div>
     
     

//       <div className="top-card-container">
//         <h2>Top Card</h2>
//         {topCard && <Card suit={topCard.Suit} value={topCard.Value} onPlay={() => {}} />}
//       </div>

//       <div className="hand-container">
//         {hand.map((card, index) => (
//           <Card key={index} suit={card.Suit} value={card.Value} onPlay={() => playCard(index)} />
//         ))}
//       </div>

//       <button className="draw-button" onClick={drawCard} disabled={hand.length >= 10}>
//         Draw Card ({hand.length}/10)
//       </button>
//     </div>
//   );
// };

// export default CrazyCardGame;


