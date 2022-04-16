import React from "react";
import { ClauseInput } from "./ClauseInput";

export class Clauses extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            clauses: [],
            model: {},
            satisfiable: false
        };
        this.handleAddClause = this.handleAddClause.bind(this);
    }

    handleAddClause(clause) {
        const clauses = this.state.clauses.slice();
        clauses.push(clause);
        this.setState({clauses: clauses});

        if(window.solveFormula) {
            const formula = clauses.join(";");
            const solution = window.solveFormula(formula);
            const satisfiable = solution["sat"];
            this.setState({satisfiable: satisfiable});
        }
    }

    render() {
        return (
            <div className="clauses">
                <ClauseInput onAddClause={this.handleAddClause} />
                <ClausesDisplay clauses={this.state.clauses} />
                <ClausesStatus isSat={this.state.satisfiable}/>
            </div>
        );
    }
}

function ClausesDisplay(props) {
    const clauses = props.clauses.map((clause, index) => {
        return <div className="clause" key={index}>{clause}</div>;
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
