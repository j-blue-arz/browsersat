import React from "react";
import { Button } from "./Button";
import { ConstraintInput } from "./ConstraintInput";
import { ConstraintsDisplay } from "./ConstraintsDisplay";

import "./Constraints.css";

export class Constraints extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            constraints: [],
            model: {},
            satisfiable: true,
            validationError: "",
            hasLoadedWasm: false,
            selectionEvaluation: undefined
        };
        this.handleAddConstraint = this.handleAddConstraint.bind(this);
        this.handleFlipLiteral = this.handleFlipLiteral.bind(this);
        this.handleClearConstraints = this.handleClearConstraints.bind(this);
        this.handleSelection = this.handleSelection.bind(this);
    }

    componentDidMount() {
        WebAssembly.instantiateStreaming(
            fetch(process.env.PUBLIC_URL + "/solver.wasm"),
            window.go.importObject
        ).then((obj) => {
            window.go.run(obj.instance);
            window.satsolver.initializeSolver();
            this.setState({ hasLoadedWasm: true });
        });
    }

    handleAddConstraint(constraint) {
        if (window.satsolver) {
            constraint = constraint.replaceAll(" ", "");
            const validationResult = window.satsolver.validateConstraint(constraint);
            if (validationResult === "VALID") {
                const constraints = this.state.constraints.slice();
                constraints.push(constraint);
                this.setState({ constraints: constraints, validationError: "" });

                window.satsolver.addConstraint(constraint);
                const satisfiable = window.satsolver.isSat();
                this.setState({ satisfiable: satisfiable });
                if (satisfiable) {
                    const model = window.satsolver.getModel();
                    this.setState({ model: model });
                }
                return true;
            } else {
                this.setState({ validationError: validationResult });
                return false;
            }
        }
    }

    handleClearConstraints() {
        this.setState({ constraints: [], model: {}, satisfiable: true });
        if (window.satsolver) {
            window.satsolver.initializeSolver();
        }
    }

    handleFlipLiteral(literal) {
        if (window.satsolver) {
            const possible = window.satsolver.flipLiteral(literal);
            if (possible) {
                const model = window.satsolver.getModel();
                this.setState({ model: model });
            }
        }
    }

    handleSelection() {
        if (window.satsolver) {
            let selection = "";

            if (window.getSelection) {
                selection = window.getSelection() + "";
            }

            if (selection !== "") {
                const evaluation = window.satsolver.evaluate(selection);
                if (evaluation !== "INVALID") {
                    this.setState({ selectionEvaluation: evaluation });
                } else {
                    this.setState({ selectionEvaluation: undefined });
                }
            }
            
        }
    }

    render() {
        return (
            <div className="constraints">
                <div className="constraints__interaction">
                    <ConstraintInput
                        onAddConstraint={this.handleAddConstraint}
                        disabled={!this.state.hasLoadedWasm}
                    />
                    <Button
                        label="Clear"
                        onClick={this.handleClearConstraints}
                        disabled={!this.state.hasLoadedWasm}
                    />
                </div>
                <ConstraintsDisplay
                    constraints={this.state.constraints}
                    model={this.state.model}
                    onFlipLiteral={this.handleFlipLiteral}
                    onSelectConstraint={this.handleSelection}
                />
                <SatStatus
                    isSat={this.state.satisfiable}
                    validationError={this.state.validationError}
                    hasLoadedWasm={this.state.hasLoadedWasm}
                    selectionEvaluation={this.state.selectionEvaluation}
                />
            </div>
        );
    }
}

function SatStatus(props) {
    if (!props.hasLoadedWasm) {
        return (
            <div className="constraints__status constraints__status--info">
                Waiting for solver to be ready...
            </div>
        );
    } else if (props.validationError !== "") {
        return (
            <div className="constraints__status constraints__status--error">
                {props.validationError}
            </div>
        );
    } else if (props.selectionEvaluation !== undefined) {
        const value = props.selectionEvaluation ? <span class="literal--true">TRUE</span> : <span class="literal--false">FALSE</span> 
        return <div className="constraints__status constraints__status--info">Selection evaluates to {value}</div>;
    } else if (props.isSat) {
        return <div className="constraints__status literal--true">SAT</div>;
    } else {
        return <div className="constraints__status literal--false">UNSAT</div>;
    }
}
