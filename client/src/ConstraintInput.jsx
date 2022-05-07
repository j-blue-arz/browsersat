import React from "react";

import "./ConstraintInput.css";

export class ConstraintInput extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            inputString: "",
        };

        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleChange = this.handleChange.bind(this);
        this.handleKeyDown = this.handleKeyDown.bind(this);
    }

    handleChange(event) {
        this.setState({ inputString: event.target.value });
    }

    handleKeyDown(event) {
        if(event.key === "Enter") {
            this.handleSubmit();
        }
    }

    handleSubmit(event) {
        if(this.state.inputString.trim !== "") {
            this.props.onAddConstraint(this.state.inputString);
        }
    }

    render() {
        return (
            <div className="constraint-input">
                <input type="text" value={this.state.inputString} onChange={this.handleChange} onKeyDown={this.handleKeyDown} />
                <button type="submit" name="add" className="constraint-input__button" onClick={this.handleSubmit}>
                    Add
                </button>
            </div>
        );
    }
}
