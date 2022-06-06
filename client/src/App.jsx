import React from "react";
import { Constraints } from "./Constraints.jsx";
import "./App.css";

function App() {
    return (
        <div className="app__wrapper">
            <div className="app__content">
                <div class="information">
                    <h4>Browsersat</h4>
                    <p>
                        Browsersat is an{" "}
                        <a href="https://github.com/j-blue-arz/browsersat">open-source</a>{" "}
                        client-side SAT-Solver.
                    </p>
                    <p>
                        Enter boolean formula(s). A satisfiable assignment will show the literals as{" "}
                        <span class="literal--true">true</span> or{" "}
                        <span class="literal--false">false</span>.{" "}
                    </p>
                    <p>
                        Clicking on a literal forces its assignment to be flipped, if this allows
                        for a satisfiable assignment. The solver will find an assignment which
                        minimizes the number of flipped literals.
                    </p>
                    <p>
                        The grammar is currently defined by the{" "}
                        <a href="https://github.com/crillab/gophersat">gophersat project</a>,
                        described{" "}
                        <a href="https://github.com/crillab/gophersat/blob/master/bf/parser.go">
                            here
                        </a>
                        . E.g. "^" is the unary negation operator, "-&gt;" is an implication.
                    </p>
                </div>
                <Constraints />
            </div>
        </div>
    );
}

export default App;
