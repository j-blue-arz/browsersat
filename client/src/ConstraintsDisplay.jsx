import React from "react";
import { Literal } from "./Literal";
import "./ConstraintsDisplay.css";

export class ConstraintsDisplay extends React.Component {

    render() {
        const constraints = this.props.constraints.map((constraint, index) => {
            const tokens = constraint.split(/([a-zA-Z0-9_]+)/);
            const display = tokens
                .filter((token) => token.length > 0)
                .map((token, index) => {
                    if (token in this.props.model) {
                        return (
                            <Literal
                                token={token}
                                key={index}
                                value={this.props.model[token]}
                                onClick={() => this.props.onFlipLiteral(token)}
                            />
                        );
                    } else {
                        return <React.Fragment key={index}>{token}</React.Fragment>;
                    }
                });

            return (
                <div className="constraint" key={index} onMouseUp={this.props.onSelectConstraint}>
                    {display}
                </div>
            );
        });
        return <div className="constraints__display">{constraints}</div>;
    }
}

ConstraintsDisplay.defaultProps = {
    model: {},
};
