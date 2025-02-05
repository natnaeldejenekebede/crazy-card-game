import React from "react";
import "./Card.css";

const suitSymbols = {
  Hearts: "♥",
  Diamonds: "♦",
  Clubs: "♣",
  Spades: "♠",
};

const Card = ({ suit, value, onPlay }) => {
  const isRed = suit === "Hearts" || suit === "Diamonds";

  return (
    <div className={`card ${isRed ? "red" : "black"}`}>
      <div>{value}</div>
      <div className="card-symbol">{suitSymbols[suit]}</div>
      <button className="card-button" onClick={onPlay}>
        Play
      </button>
    </div>
  );
};

export default Card;
