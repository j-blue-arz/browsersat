import React from "react";
import { Button } from "./Button.jsx";

import "./ConstraintInput.css";

export class ConstraintInput extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            inputString: "a->(b|c)",
        };

        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleChange = this.handleChange.bind(this);
        this.handleKeyDown = this.handleKeyDown.bind(this);
    }

    handleChange(event) {
        this.setState({ inputString: event.target.value });
    }

    handleKeyDown(event) {
        if (event.key === "Enter") {
            this.handleSubmit();
        }
    }

    handleSubmit(event) {
        if (this.state.inputString.trim !== "") {
            this.props.onAddConstraint(this.state.inputString);
        }
    }

    render() {
        return (
            <div className="constraint-input">
                <label for="constraint-input">Enter constraint:</label>
                <input
                    id="constraint-input"
                    type="text"
                    className="constraint-input__text"
                    value={this.state.inputString}
                    onChange={this.handleChange}
                    onKeyDown={this.handleKeyDown}
                />
                <Button label="Add" onClick={this.handleSubmit} />
            </div>
        );
    }
}
