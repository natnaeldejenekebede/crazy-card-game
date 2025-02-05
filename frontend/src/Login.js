// import React, { useState } from 'react';
// //import './Auth.module.css';  // Import the Auth.css file
// import { useNavigate } from "react-router-dom";

// const Login = ({ onSwitch }) => {
//   const [username, setUsername] = useState('');
//   const [password, setPassword] = useState('');
//   const [message, setMessage] = useState('');
//   const [token, setToken] = useState('');
//   const navigate = useNavigate(); 
//   const handleSubmit = async (e) => {
//     e.preventDefault();
//     console.log("here")
//     const response = await fetch('http://localhost:8083/login', {
//       method: 'POST',
//       headers: {
//         'Content-Type': 'application/json',
//       },
//       body: JSON.stringify({ username, password }),
//     });

//     const data = await response.json();
//     if (response.ok) {
//       setMessage('Login successful');
//       setToken(data.token);  // Assuming the token is returned in the response
//       setUsername('');
//       setPassword('');
//       navigate("/roomselection");
//     } else {
//         console.log("abjhvjv")
//       setMessage(data.error || 'Error logging in');
//     }
//   };

//   return (
//     <div className="auth-container">
//       <h2>Login</h2>
//       <form onSubmit={handleSubmit}>
//         <input
//           type="text"
//           placeholder="Username"
//           value={username}
//           onChange={(e) => setUsername(e.target.value)}
//           required
//         />
//         <input
//           type="password"
//           placeholder="Password"
//           value={password}
//           onChange={(e) => setPassword(e.target.value)}
//           required
//         />
//         <button type="submit">Login</button>
//       </form>
//       {message && <p className="message">{message}</p>}
//       {token && <p className="message">Your token: {token}</p>} {/* Optional, just to show the token */}
//       <p>
//         Don't have an account?{' '}
//         <button onClick={onSwitch}>Sign up here</button>
//       </p>
//     </div>
//   );
// };

// export default Login;



// import React, { useState } from 'react';
// //import './Auth.module.css';  // Import the Auth.css file
// import { useNavigate } from "react-router-dom";

// const Login = ({ onSwitch }) => {
//   const [username, setUsername] = useState('');
//   const [password, setPassword] = useState('');
//   const [message, setMessage] = useState('');
//   const [token, setToken] = useState('');
//   const navigate = useNavigate(); 
//   const handleSubmit = async (e) => {
//     e.preventDefault();
//     console.log("here")
//     const response = await fetch('http://localhost:8083/login', {
//       method: 'POST',
//       headers: {
//         'Content-Type': 'application/json',
//       },
//       body: JSON.stringify({ username, password }),
//     });

//     const data = await response.json();
//     console.log(data.token)
//     // localStorage.setItem("token", data.token);
//     sessionStorage.setItem("token", data.token); 
//     if (response.ok) {
//       setMessage('Login successful');
//       setToken(data.token);  // Assuming the token is returned in the response
//       // localStorage.setItem("username", username);
//       sessionStorage.setItem("username", username); 
//       setUsername('');
//       setPassword('');
//       navigate("/roomselection");
//     } else {
//         console.log("abjhvjv")
//       setMessage(data.error || 'Error logging in');
//     }
//   };

//   return (
//     <div className="auth-container">
//       <h2>Login</h2>
//       <form onSubmit={handleSubmit}>
//         <input
//           type="text"
//           placeholder="Username"
//           value={username}
//           onChange={(e) => setUsername(e.target.value)}
//           required
//         />
//         <input
//           type="password"
//           placeholder="Password"
//           value={password}
//           onChange={(e) => setPassword(e.target.value)}
//           required
//         />
//         <button type="submit">Login</button>
//       </form>
//       {message && <p className="message">{message}</p>}
//       {token && <p className="message">Your token: {token}</p>} {/* Optional, just to show the token */}
//       <p>
//         Don't have an account?{' '}
//         <button onClick={onSwitch}>Sign up here</button>
//       </p>
//     </div>
//   );
// };

// export default Login;

import React, { useState } from 'react';
import { useNavigate } from "react-router-dom";
import './Login.css';  // Import the new CSS file

const Login = ({ onSwitch }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [token, setToken] = useState('');
  const navigate = useNavigate(); 

  const handleSubmit = async (e) => {
    e.preventDefault();
    console.log("Logging in...");
    
    const response = await fetch('http://localhost:8083/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    const data = await response.json();
    console.log(data.token);

    sessionStorage.setItem("token", data.token);
    
    if (response.ok) {
      setMessage('Login successful');
      setToken(data.token);  
      sessionStorage.setItem("username", username); 
      setUsername('');
      setPassword('');
      navigate("/roomselection");
    } else {
      console.log("Login error");
      setMessage(data.error || 'Error logging in');
    }
  };

  return (
    <div className="login-container">
      <h2 className="login-title">Login</h2>
      <form className="login-form" onSubmit={handleSubmit}>
        <input
          className="login-input"
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          className="login-input"
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button className="login-button" type="submit">Login</button>
      </form>
      {message && <p className="login-message">{message}</p>}
      {token && <p className="login-token">Your token: {token}</p>} {/* Optional, just to show the token */}
      <p className="login-switch">
        Don't have an account?{' '}
        <button className="login-switch-btn" onClick={onSwitch}>Sign up here</button>
      </p>
    </div>
  );
};

export default Login;
