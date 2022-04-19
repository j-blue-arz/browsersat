import React from "react";
import { ConstraintInput } from "./ConstraintInput";
import { ConstraintsDisplay } from "./ConstraintsDisplay";

export class Constraints extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            constraints: [],
            model: {},
            satisfiable: false
        };
        this.handleAddConstraint = this.handleAddConstraint.bind(this);
    }

    handleAddConstraint(constraint) {
        const constraints = this.state.constraints.slice();
        constraints.push(constraint);
        this.setState({constraints: constraints});

        if(window.solveFormula) {
            const formula = constraints.join(";");
            const solution = window.solveFormula(formula);
            const satisfiable = solution["sat"];
            this.setState({satisfiable: satisfiable});
            if(satisfiable) {
                this.setState({model: solution["model"]})
            }
        }
    }

    render() {
        return (
            <div className="constraints">
                <ConstraintInput onAddConstraint={this.handleAddConstraint} />
                <ConstraintsDisplay constraints={this.state.constraints} model={this.state.model} />
                <SatStatus isSat={this.state.satisfiable}/>
            </div>
        );
    }
}


function SatStatus(props) {
    if(props.isSat) {
        return <div className="constraints__status constraints__status--sat">SAT</div>;
    } else {
        return <div className="constraints__status constraints__status--unsat">UNSAT</div>;
    }
    
}
