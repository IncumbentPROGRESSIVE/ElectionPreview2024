import React, { useState, useEffect } from "react";
import axios from "axios";
import "./App.css";

const SignupComponent = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loggedInUser, setLoggedInUser] = useState(null);
  const [showUserBox, setShowUserBox] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      axios
        .get("http://localhost:8080/profile", {
          headers: { Authorization: `Bearer ${token}` },
        })
        .then((response) => {
          setLoggedInUser(response.data.email);
        })
        .catch((error) => {
          console.error("Error verifying token:", error);
        });
    }
  }, []);

  const handleRegister = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post("http://localhost:8080/register", {
        email,
        password,
      });
      alert("User registered successfully");
      setLoggedInUser(email);
      localStorage.setItem("token", response.data.token);
    } catch (error) {
      console.error("Error registering user:", error.response || error);
      alert("Error registering user");
    }
  };

  const handleLogin = async (e) => {
    e.preventDefault();
    try {
      const response = await axios.post("http://localhost:8080/login", {
        email,
        password,
      });
      const { token } = response.data;
      localStorage.setItem("token", token);
      setLoggedInUser(email);
    } catch (error) {
      console.error("Error logging in user:", error.response || error);
      alert("Error logging in user");
    }
  };

  const handleLogout = () => {
    setLoggedInUser(null);
    localStorage.removeItem("token");
  };

  return (
    <div className="signup-container">
      {!loggedInUser ? (
        <form onSubmit={handleLogin} className="signup-form">
          <div className="form-group">
            <label htmlFor="email">Email Address:</label>
            <input
              type="email"
              id="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="password">Password:</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="save-info-button">
            Sign In
          </button>
          <button
            type="button"
            onClick={handleRegister}
            className="save-info-button"
            style={{ marginTop: "10px" }}
          >
            Sign Up
          </button>
        </form>
      ) : (
        <div>
          <button onClick={handleLogout} className="save-info-button">
            Logout
          </button>
          <button
            onClick={() => setShowUserBox(!showUserBox)}
            className="save-info-button"
          >
            Show User
          </button>
          {showUserBox && (
            <div className="user-box">
              <p>Logged in as: {loggedInUser}</p>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default SignupComponent;
