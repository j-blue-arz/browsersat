import React from "react";
import { render, unmountComponentAtNode } from "react-dom";
import { act } from "react-dom/test-utils";

import { ConstraintsDisplay } from "./ConstraintsDisplay";

let container = null;
let constraints = null;
let model = {};

beforeEach(() => {
    // setup a DOM element as a render target
    container = document.createElement("div");
    document.body.appendChild(container);
});

describe("ConstraintsDisplay", () => {
    it("renders constraint as text", () => {
        givenConstraints(["a -> b | c"]);

        act(() => {
            render(<ConstraintsDisplay constraints={constraints} />, container);
        });

        expect(container.textContent).toBe("a -> b | c");
        expect(container.getElementsByClassName("constraint").length).toBe(1);
    });

    it("surrounds literals in model with a span", () => {
        givenConstraints(["abc -> edf | crg"]);
        givenModel({ abc: true, edf: true, crg: true });

        act(() => {
            render(<ConstraintsDisplay constraints={constraints} model={model} />, container);
        });

        const constraintDiv = container.getElementsByClassName("constraint")[0];
        expect(constraintDiv.getElementsByTagName("span").length).toBe(3);
        expectSpanWithContent(constraintDiv.childNodes[0], "abc");
        expectSpanWithContent(constraintDiv.childNodes[2], "edf");
        expectSpanWithContent(constraintDiv.childNodes[4], "crg");
    });

    it("marks literals in model according to their value", () => {
        givenConstraints(["a -> b | c"]);
        givenModel({ a: false, b: true, c: true });

        act(() => {
            render(<ConstraintsDisplay constraints={constraints} model={model} />, container);
        });

        expect(findOneSpanContaining("a").className).toContain("literal--false");
        expect(findOneSpanContaining("b").className).toContain("literal--true");
        expect(findOneSpanContaining("c").className).toContain("literal--true");
        
    });
});

function givenConstraints(given) {
    constraints = given;
}

function givenModel(given) {
    model = given;
}

function expectSpanWithContent(element, content) {
    expect(element.tagName).toBe("SPAN");
    expect(element.innerHTML).toBe(content);
}

function findOneSpanContaining(content) {
    const spans = Array.from(container.getElementsByTagName("span"));
    const matches = spans.filter(element => element.innerHTML === content);
    expect(matches.length).toBe(1);
    return matches[0];
}

afterEach(() => {
    // cleanup on exiting
    unmountComponentAtNode(container);
    container.remove();
    container = null;
});
