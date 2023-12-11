import ReactMarkdown from "react-markdown";
import { Message } from "./data";

export interface Props {
  message: Message;
  onContinueClicked: (() => void) | null;
  onDeleteClicked: (() => void) | null;
}

export function ResponseMessage({message, onContinueClicked, onDeleteClicked}: Props): JSX.Element {
  const iconName = message.role === "Assistant" ? "fa-robot" : "fa-user";
  return <div className="response-message">
    <div className="response-message-gutter"><i className={"fas " + iconName}></i></div>
    <div className="response-message-text">
      <div className="response-message-controls">
      {
        onContinueClicked != null &&
          <button className="compact small" onClick={onContinueClicked}>Continue</button>
      }
      {
        onDeleteClicked != null &&
          <button className="microtool danger" onClick={onDeleteClicked}><i className="fa fa-times"></i></button>
      }
      </div>
      <ReactMarkdown children={message.text} /><br />
    </div>
  </div>;
}
