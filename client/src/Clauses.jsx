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
        return (
            <div class="clauses">
                <ClauseInput onAddClause={this.handleAddClause} />
                <ClausesDisplay clauses={this.state.clauses} />
                <ClausesStatus />
            </div>
        );
    }
}

function ClausesDisplay(props) {
    console.log("display" + props.clauses);
    const clauses = props.clauses.map((clause) => {
        return <div class="clause">{clause}</div>;
    });
    return <div class="clauses__display">{clauses}</div>;
}

function ClausesStatus() {
    return <div class="clauses__status">UNSAT</div>;
}
