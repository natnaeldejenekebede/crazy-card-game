// import logo from './logo.svg';
// import './App.css';


// function App() {
//   return (
//     <div className="App">
//       <header className="App-header">
//         <img src={logo} className="App-logo" alt="logo" />
//         <p>
//           Edit <code>src/App.js</code> and save to reload.
//         </p>
//         <a
//           className="App-link"
//           href="https://reactjs.org"
//           target="_blank"
//           rel="noopener noreferrer"
//         >
//           Learn React
//         </a>
//       </header>
//     </div>
//   );
// }

// export default App;
import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import CrazyCardGame from "./CrazyCardGame"
import AuthPage from './AuthPage';
import RoomSelection from './RoomSelection'; // Import your next page

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<AuthPage />} />
        <Route path="/roomselection" element={<RoomSelection />} />
        <Route path="/card/:suit" element={<CrazyCardGame />} />
      </Routes>
    </Router>
  );
};

export default App;
