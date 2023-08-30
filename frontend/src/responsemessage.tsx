import ReactMarkdown from "react-markdown";
import { Message } from "./data";

export interface Props {
  message: Message;
  onContinueClicked: (() => void) | null;
}

export function ResponseMessage({message, onContinueClicked}: Props): JSX.Element {
  const iconName = message.role === "Assistant" ? "fa-robot" : "fa-user";
  return <div className="response-message">
    <div className="response-message-gutter"><i className={"fas " + iconName}></i></div>
    <div className="response-message-text">
      <ReactMarkdown children={message.text} /><br />
      {
        onContinueClicked != null &&
          <button className="compact small" onClick={onContinueClicked}>Continue</button>

      }
    </div>
  </div>;
}
