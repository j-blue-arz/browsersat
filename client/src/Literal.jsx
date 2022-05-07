import React from "react";
import "./Literal.css";

export class Literal extends React.Component {
    render() {
        let valueClass = "";
        if (this.props.value === true) {
            valueClass = "literal--true";
        } else if (this.props.value === false) {
            valueClass = "literal--false";
        }
        return <span className={`literal ${valueClass}`} onClick={this.props.onClick}>{this.props.token}</span>;
    }
}
