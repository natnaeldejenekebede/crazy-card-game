import React, { useState } from 'react';
import Signup from './Signup';
import Login from './Login';

const AuthPage = () => {
  const [isSignup, setIsSignup] = useState(true);

  const toggleForm = () => {
    setIsSignup(!isSignup);
  };

  return (
    <div>
      {/* <h1>{isSignup ? 'Create Account' : 'Login to Your Account'}</h1> */}
      {isSignup ? (
        <Signup onSwitch={toggleForm} />
      ) : (
        <Login onSwitch={toggleForm} />
      )}
    </div>
  );
};

export default AuthPage;
