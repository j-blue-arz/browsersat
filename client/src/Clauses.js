import React from "react";

export function Clauses() {
    return (
        <div class="clauses">
            <ClausesInputarea />
            <ClausesDisplay />
            <ClausesStatus />
        </div>
    );
}

class ClausesInputarea extends React.Component {
    render() {
        return (
            <div class="clauses__inputarea">
                <input id="clause-field" type="text" />
                <button id="add-clause" name="add">
                    Add
                </button>
            </div>
        );
    }
}

function ClausesDisplay() {
    return (
        <div class="clauses__display">
            <div class="clause">a</div>
        </div>
    );
}

function ClausesStatus() {
    return (
        <div class="clauses__status">UNSAT</div>
    )
}
