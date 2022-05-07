import React from "react";
import { ConstraintInput } from "./ConstraintInput";
import { ConstraintsDisplay } from "./ConstraintsDisplay";

import "./Constraints.css";

export class Constraints extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            constraints: [],
            model: {},
            satisfiable: false,
        };
        this.handleAddConstraint = this.handleAddConstraint.bind(this);
        this.handleFlipLiteral = this.handleFlipLiteral.bind(this);
    }

    handleAddConstraint(constraint) {
        if (window.satsolver) {
            constraint = constraint.replaceAll(" ", "");
            window.satsolver.addConstraint(constraint);

            const constraints = this.state.constraints.slice();
            constraints.push(constraint);
            this.setState({ constraints: constraints });

            const satisfiable = window.satsolver.isSat()
            this.setState({ satisfiable: satisfiable });
            if (satisfiable) {
                const model = window.satsolver.getModel();
                this.setState({ model: model });
            }
        }
    }

    handleFlipLiteral(literal) {
        if (window.satsolver) {
            const possible = window.satsolver.flipLiteral(literal);
            if(possible) {
                const model = window.satsolver.getModel();
                this.setState({ model: model });
            }
        }
    }

    render() {
        return (
            <div className="constraints">
                <ConstraintInput onAddConstraint={this.handleAddConstraint} />
                <ConstraintsDisplay constraints={this.state.constraints} model={this.state.model} onFlipLiteral={this.handleFlipLiteral} />
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
