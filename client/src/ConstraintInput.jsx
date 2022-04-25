import React from "react";

import "./ConstraintInput.css";

export class ConstraintInput extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            inputString: "",
        };

        this.handleOnClick = this.handleOnClick.bind(this);
        this.handleChange = this.handleChange.bind(this);
    }

    handleOnClick(event) {
        this.props.onAddConstraint(this.state.inputString);
    }

    handleChange(event) {
        this.setState({ inputString: event.target.value });
    }

    render() {
        return (
            <div className="constraint-input">
                <input type="text" value={this.state.inputString} onChange={this.handleChange} />
                <button type="submit" name="add" className="constraint-input__button" onClick={this.handleOnClick}>
                    Add
                </button>
            </div>
        );
    }
}
