import ReactMarkdown from "react-markdown";
import { Message } from "./data";
import { FileAttachmentsList } from "./fileattachmentslist";

export interface Props {
  sessionId: string;
  message: Message;
  onContinueClicked: (() => void) | null;
  onDeleteClicked: (() => void) | null;
}

export function ResponseMessage({sessionId, message, onContinueClicked, onDeleteClicked}: Props): JSX.Element {
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
      {
        message.attachedFiles != null && message.attachedFiles.length != 0 &&
        <FileAttachmentsList
          sessionId={sessionId}
          attachedFiles={message.attachedFiles}
        />
      }
    </div>
  </div>;
}
