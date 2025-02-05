import React from "react";
import "./OpponentCards.css";

const OpponentCards = ({ opponents}) => {
    return (
      <div className="opponents-container">
        <h2>Opponents</h2>
        <div className="opponents-list" style={{ display: 'flex', flexDirection: 'row', gap: '5px' }}>
          {[...Array(opponents.num)].map((_, index) => (
            <div key={index} className="card-back" />
          ))}
        </div>
        <p>Cards Left: {opponents.num}</p>
      </div>
    );
  };
  
  export default OpponentCards;

// const OpponentCards = ({ opponents }) => {
//   return (
//     <div className="opponents-container">
//       <h2>Opponents</h2>
//       <div className="opponents-list" style={{ display: 'flex', flexDirection: 'row', gap: '5px' }}>
//         {opponents.map((opponent, index) => (
//           <div key={index} className="card-holder">
//             {/* Render the number of cards each opponent has */}
//             {[...Array(opponent.num)].map((_, cardIndex) => (
//               <div key={cardIndex} className="card-back" style={{ width: '50px', height: '75px', backgroundColor: 'gray' }} />
//             ))}
//           </div>
//         ))}
//       </div>
//       <p>Cards Left: {opponents.length}</p>
//     </div>
//   );
// };

// export default OpponentCards;

