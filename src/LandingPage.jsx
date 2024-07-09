// LandingPage.jsx
import React from "react";
import "./App.css"; // Import the App.css file
//import Flair from "./flair.jsx";
import SignupComponent from "./saveinfo";

function LandingPage() {
  return (
    <div className="landing-page">
      {/*<h1 className="custom-heading">Election Preview 2024</h1>{" "}

      <Flair /> */}
      <SignupComponent />
    </div>
  );
}

export default LandingPage;
