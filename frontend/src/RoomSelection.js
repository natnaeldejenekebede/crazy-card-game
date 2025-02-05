// import React, { useState, useEffect } from "react";
// import { useNavigate } from "react-router-dom";
// //import "./RoomSelection.css"; // Importing regular CSS file

// function RoomSelection() {
//   const navigate = useNavigate();
//   const [rooms, setRooms] = useState([]);
//   const [openingRoom, setOpeningRoom] = useState(null);

//   useEffect(() => {
//     // Simulating fetching room data, with random availability (open/closed)
//     const fetchedRooms = Array.from({ length: 4 }, (_, i) => ({
//       number: i + 1,
//       suit: ["hearts", "diamonds", "clubs", "spades"][i],
//       symbol: ["♥", "♦", "♣", "♠"][i],
//       isOpen: Math.random() > 0.5, // Random open/closed rooms
//     }));
//     setRooms(fetchedRooms);
//   }, []);

//   const enterRoom = (room) => {
//     if (!room.isOpen) {
//       alert("Room is full! Choose another.");
//       return;
//     }

//     setOpeningRoom(room.number);
//     setTimeout(() => {
//       navigate(`/crazycardgame/${room.number}`);
//     }, 500);
//   };

//   return (
//     <div className="roomContainer">
//       <h1>Select a Room</h1>
//       <div className="roomList">
//         {rooms.map((room) => (
//           <div
//             key={room.number}
//             className={`room 
//                         ${room.isOpen ? "openRoom" : "closedRoom"} 
//                         ${openingRoom === room.number ? "opening" : ""} 
//                         ${room.suit}`} // Applying the suit as a class
//             onClick={() => enterRoom(room)}
//           >
//             <span className="symbol">{room.symbol}</span>
//             <span className="roomText">{room.suit.charAt(0).toUpperCase() + room.suit.slice(1)} Room</span>
//           </div>
//         ))}
//       </div>
//     </div>
//   );
  
// }

// export default RoomSelection;




// import React, { useState } from "react";
// import { useNavigate } from "react-router-dom";
// import axios from "axios";

// const authServerURL = "http://localhost:8083"; // Replace with actual auth server URL

// function RoomSelection() {
//   const navigate = useNavigate();
//   const rooms = [
//     { suit: "hearts", symbol: "♥" },
//     { suit: "diamonds", symbol: "♦" },
//     { suit: "clubs", symbol: "♣" },
//     { suit: "spades", symbol: "♠" },
//   ];

//   const enterRoom = async (room) => {
//     // const token = localStorage.getItem("token"); // Retrieve token
//     const token = sessionStorage.getItem("token");
//     if (!token) {
//       alert("Authentication error: No token found. Please log in again.");
//       return;
//     }

//     try {
//       const response = await axios.get(`${authServerURL}/roomselection?suit=${room.suit}`, {
//         headers: {
//           Authorization: `Bearer ${token}`, // Properly formatted token
//         },
//       });
//       localStorage.setItem("ip", response.data.ip);

//       console.log("Server response:", response.data);

//       if (response.data.message === "good") {
//         navigate(`/card/${room.suit}`); // Navigate to the card page
//       } else {
//         alert("Room is full or unavailable.");
//       }
//     } catch (error) {
//       console.error("Error checking room:", error);
//       alert("Failed to join room. Try again.");
//     }
//   };

//   return (
//     <div className="roomContainer">
//       <h1>Select a Room</h1>
//       <div className="roomList">
//         {rooms.map((room) => (
//           <div
//             key={room.suit}
//             className="room"
//             onClick={() => enterRoom(room)}
//           >
//             <span className="symbol">{room.symbol}</span>
//             <span className="roomText">{room.suit.charAt(0).toUpperCase() + room.suit.slice(1)} Room</span>
//           </div>
//         ))}
//       </div>
//     </div>
//   );
// }

// export default RoomSelection;


import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import "./RoomSelection.css"; // Import the CSS file

const authServerURL = "http://localhost:8083"; // Replace with actual auth server URL

function RoomSelection() {
  const navigate = useNavigate();
  const rooms = [
    { suit: "hearts", symbol: "♥" },
    { suit: "diamonds", symbol: "♦" },
    { suit: "clubs", symbol: "♣" },
    { suit: "spades", symbol: "♠" },
  ];

  const enterRoom = async (room) => {
    const token = sessionStorage.getItem("token");
    if (!token) {
      alert("Authentication error: No token found. Please log in again.");
      return;
    }

    try {
      const response = await axios.get(`${authServerURL}/roomselection?suit=${room.suit}`, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
      localStorage.setItem("ip", response.data.ip);

      if (response.data.message === "good") {
        navigate(`/card/${room.suit}`);
      } else {
        alert("Room is full or unavailable.");
      }
    } catch (error) {
      console.error("Error checking room:", error);
      alert("Failed to join room. Try again.");
    }
  };

  return (
    <div className="room-selection-container">
      <h1>Select a Room</h1>
      <div className="room-selection-list">
        {rooms.map((room) => (
          <div
            key={room.suit}
            className="room-selection-item"
            onClick={() => enterRoom(room)}
          >
            <span className="symbol">{room.symbol}</span>
            <span className="roomText">{room.suit.charAt(0).toUpperCase() + room.suit.slice(1)} Room</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export default RoomSelection;
