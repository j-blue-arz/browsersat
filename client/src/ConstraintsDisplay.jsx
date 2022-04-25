import React from "react";

import "./ConstraintsDisplay.css";

export class ConstraintsDisplay extends React.Component {
    render() {

        const constraints = this.props.constraints.map((constraint, index) => {
            const tokens = constraint.split(/([a-zA-Z0-9_]+)/);
            const display = tokens.filter(token => token.length > 0).map((token, index) => {
                if (token in this.props.model) {
                    return <Literal token={token} key={token} value={this.props.model[token]} />
                } else {
                    return <React.Fragment key={index}>{token}</React.Fragment>
                }
            });

            return (
                <div className="constraint" key={index}>
                    {display}
                </div>
            );
        });
        return <div className="constraints__display">{constraints}</div>;
    }
}

function Literal(props) {
    let valueClass = "";
    if(props.value === true) {
        valueClass = "literal--true";
    } else if(props.value === false) {
        valueClass = "literal--false";
    }
    return <span className={`literal ${valueClass}`}>{props.token}</span>
}

ConstraintsDisplay.defaultProps = {
    model: {}
}
