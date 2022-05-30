import "./Button.css";

export function Button(props) {
    return (
        <button type="submit" className="button-white" onClick={props.onClick}>
            {props.label}
        </button>
    );
}
