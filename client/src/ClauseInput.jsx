import React from "react";

export class ClauseInput extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            clauseString: "",
        };

        this.handleOnClick = this.handleOnClick.bind(this);
        this.handleChange = this.handleChange.bind(this);
    }

    handleOnClick(event) {
        this.props.onAddClause(this.state.clauseString);
    }

    handleChange(event) {
        this.setState({ clauseString: event.target.value });
    }

    render() {
        return (
            <div class="clause-input">
                <input type="text" value={this.state.clauseString} onChange={this.handleChange} />
                <button name="add" onClick={this.handleOnClick}>
                    Add
                </button>
            </div>
        );
    }
}
