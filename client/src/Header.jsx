import React from "react";

import "./Header.css";

export class Header extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            showInfo: false,
        };

        this.toggleInfo = this.toggleInfo.bind(this);
    }

    toggleInfo() {
        this.setState({ showInfo: !this.state.showInfo });
    }

    render() {
        let content;
        if (this.state.showInfo) {
            content = <LongInfo onHideInfo={this.toggleInfo}/>;
        } else {
            content = <Menu onShowInfo={this.toggleInfo}/>;
        }
        return (
            <div className="header">
                {content}
            </div>
        );
    }
}

function LongInfo(props) {
    return (
        <div className="header__information">
            <p className="header__menu">
                Browsersat | <span onClick={props.onHideInfo}>hide</span>
            </p>
            <p>
                <span className="header__projectname">Browsersat</span> is an{" "}
                <a href="https://github.com/j-blue-arz/browsersat">open-source</a> client-side
                SAT-Solver.
            </p>
            <p>
                Enter boolean formula(s). A satisfiable assignment will show the literals as{" "}
                <span className="literal--true">true</span> or{" "}
                <span className="literal--false">false</span>.{" "}
            </p>
            <p>
                Clicking on a literal forces its assignment to be flipped, if this allows for a
                satisfiable assignment. The solver will find an assignment which minimizes the
                number of flipped literals.
            </p>
            <p>
                The grammar is currently defined by the{" "}
                <a href="https://github.com/crillab/gophersat">gophersat project</a>, described{" "}
                <a href="https://pkg.go.dev/github.com/crillab/gophersat@v1.3.1/bf#Parse">here</a>.
                E.g. "^" is the unary negation operator, "-&gt;" is an implication. Uniqueness is
                currently not supported, though.
            </p>
        </div>
    );
}

function Menu(props) {
    return (
        <p className="header__menu">
            Browsersat | <span onClick={props.onShowInfo}>info</span>
        </p>
    );
}
