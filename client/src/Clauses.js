import React from "react";

export function Clauses() {
    return (
        <div class="clauses">
            <div class="clauses__inputarea">
                <input id="clause-field" type="text"/>
                <button id="add-clause" name="add">Add</button>
                is satisfiable: <output id="satisfiable" for="add"></output>
            </div>
            <div class="clauses__outputarea">
                <div class="clause">
                    a
                </div>
            </div>
        </div>
    )
}