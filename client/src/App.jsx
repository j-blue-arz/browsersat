import React from "react";
import { Constraints } from "./Constraints.jsx";
import { Header } from "./Header.jsx";
import "./App.css";

function App() {
    return (
        <div className="app__wrapper">
            <div className="app__content">
                <Header />
                <Constraints />
            </div>
        </div>
    );
}

export default App;
