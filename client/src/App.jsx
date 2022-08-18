import React from "react";
import { Constraints } from "./Constraints.jsx";
import { Header } from "./Header.jsx";
import "./App.css";

export default class App extends React.Component {
    componentDidMount() {
        WebAssembly.instantiateStreaming(fetch(process.env.PUBLIC_URL + "/solver.wasm"), window.go.importObject).then(
            (obj) => {
                window.go.run(obj.instance);
                window.satsolver.initializeSolver();
            }
        );
    }

    render() {
        return (
            <div className="app__wrapper">
                <div className="app__content">
                    <Header />
                    <Constraints />
                </div>
            </div>
        );
    }
}
