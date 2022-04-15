import React from "react";
import { ClauseInput } from "./ClauseInput";

export class Clauses extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            clauses: [],
        };
        this.handleAddClause = this.handleAddClause.bind(this);
    }

    handleAddClause(clause) {
        const clauses = this.state.clauses.slice();
        clauses.push(clause);
        this.setState({clauses: clauses});
    }

    render() {
        const formula = this.state.clauses.join(";");
        let result = false;
        if(window.solveFormula) {
            result = window.solveFormula(formula);
        }
        return (
            <div className="clauses">
                <ClauseInput onAddClause={this.handleAddClause} />
                <ClausesDisplay clauses={this.state.clauses} />
                <ClausesStatus isSat={result}/>
            </div>
        );
    }
}

function ClausesDisplay(props) {
    const clauses = props.clauses.map((clause) => {
        return <div className="clause">{clause}</div>;
    });
    return <div className="clauses__display">{clauses}</div>;
}

function ClausesStatus(props) {
    if(props.isSat) {
        return <div className="clauses__status clauses__status--sat">SAT</div>;
    } else {
        return <div className="clauses__status clauses__status--unsat">UNSAT</div>;
    }
    
}
