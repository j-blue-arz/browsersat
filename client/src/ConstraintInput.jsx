import React from "react";
import { Button } from "./Button.jsx";

import "./ConstraintInput.css";

export class ConstraintInput extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            inputString: "a->(b|c)",
            invalidInput: false
        };

        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleChange = this.handleChange.bind(this);
        this.handleKeyDown = this.handleKeyDown.bind(this);
    }

    handleChange(event) {
        this.setState({ inputString: event.target.value, invalidInput: false });
    }

    handleKeyDown(event) {
        if (event.key === "Enter") {
            this.handleSubmit();
        }
    }

    handleSubmit(event) {
        if (this.state.inputString.trim() !== "") {
            const success = this.props.onAddConstraint(this.state.inputString);
            if(success){
                this.setState({ inputString: "", invalidInput: false });
            } else {
                this.setState({ invalidInput: true });
            }
        }
    }

    render() {
        return (
            <div className="constraint-input">
                <label className="constraint-input__label">
                    Enter constraint:
                    <input
                        type="text"
                        className={`constraint-input__field ${this.state.invalidInput ? "constraint-input__field--invalid" : ""}`}
                        value={this.state.inputString}
                        onChange={this.handleChange}
                        onKeyDown={this.handleKeyDown}
                    />
                </label>
                <Button label="Add" onClick={this.handleSubmit} disabled={this.props.disabled}/>
            </div>
        );
    }
}
