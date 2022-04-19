import React from "react";
import { render, unmountComponentAtNode } from "react-dom";
import { act } from "react-dom/test-utils";

import { ConstraintsDisplay } from "./ConstraintsDisplay";

let container = null;
beforeEach(() => {
    // setup a DOM element as a render target
    container = document.createElement("div");
    document.body.appendChild(container);
});

describe("ConstraintsDisplay", () => {
    it("renders constraint as text", () => {
        const constraints = ["a -> b | c"];
        act(() => {
            render(<ConstraintsDisplay constraints={constraints} />, container);
        });
        expect(container.textContent).toBe("a -> b | c");
    });

    it("surrounds literals in model with a span", () => {
        const constraints = ["a -> b | c"];
        const model = {"a": true, "b": true, "c": true}

        act(() => {
            render(<ConstraintsDisplay constraints={constraints} model={model} />, container);
        });

        expect(container.getElementsByTagName('span').length).toBe(3);
    })


});

afterEach(() => {
    // cleanup on exiting
    unmountComponentAtNode(container);
    container.remove();
    container = null;
});
