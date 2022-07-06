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
        };
        this.handleAddConstraint = this.handleAddConstraint.bind(this);
        this.handleFlipLiteral = this.handleFlipLiteral.bind(this);
        this.handleClearConstraints = this.handleClearConstraints.bind(this);
    }

    handleAddConstraint(constraint) {
        if (window.satsolver) {
            constraint = constraint.replaceAll(" ", "");
            const result = window.satsolver.addConstraint(constraint);
            if (result === true) {
                const constraints = this.state.constraints.slice();
                constraints.push(constraint);
                this.setState({ constraints: constraints });

                const satisfiable = window.satsolver.isSat();
                this.setState({ satisfiable: satisfiable });
                if (satisfiable) {
                    const model = window.satsolver.getModel();
                    this.setState({ model: model });
                }
                return true;
            } else {
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

    render() {
        return (
            <div className="constraints">
                <div className="constraints__interaction">
                    <ConstraintInput onAddConstraint={this.handleAddConstraint} />
                    <Button label="Clear" onClick={this.handleClearConstraints} />
                </div>
                <ConstraintsDisplay
                    constraints={this.state.constraints}
                    model={this.state.model}
                    onFlipLiteral={this.handleFlipLiteral}
                />
                <SatStatus isSat={this.state.satisfiable} />
            </div>
        );
    }
}

function SatStatus(props) {
    if (props.isSat) {
        return <div className="constraints__status constraints__status--sat">SAT</div>;
    } else {
        return <div className="constraints__status constraints__status--unsat">UNSAT</div>;
    }
}
