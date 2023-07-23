import { Response } from "./data";
import { ChangeEvent, KeyboardEventHandler, useState } from "react";
import classNames from "classnames";
import { ResponseMessage } from "./responsemessage";
import TextareaAutosize from "react-textarea-autosize";

export interface Props {
  response: Response;
  onDeleteClicked: (responseId: string) => void;
  onReply: (replyText: string) => void;
}

export function ResponseEditor({response, onDeleteClicked, onReply}: Props): JSX.Element {
  const [isPromptOpen, setIsPromptOpen] = useState<boolean>(false);
  const [reply, setReply] = useState<string>("");

  const onPromptClicked = () => {
    setIsPromptOpen(!isPromptOpen);
  };

  const onReplyChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    setReply(event.target.value);
  }

  const onReplyKeyDown: KeyboardEventHandler<HTMLTextAreaElement> = (e) => {
    if (e.code === "Enter" && e.shiftKey) {
      onReply(reply);
      setReply("");
      e.preventDefault();
    }
  };

  const onReplyClicked = () => {
    onReply(reply);
    setReply("");
  };

  return <div className="card char-width-20">
    <h3>Response</h3>
    <div className="controls">
      { response.status === "Pending" && <span className="badge warning">Pending</span>}
      { response.status === "Running" && <span className="badge success">Running</span>}
      { response.status === "Error" && <span className="badge danger">Error</span>}
      <button className="microtool danger" onClick={() => onDeleteClicked(response.id)}><i className="fa fa-times"></i></button>
    </div>
    {response.messages.length !==0 &&
      <h4 className="prompt-header" onClick={onPromptClicked}><i className={classNames({"fa": true, "fa-chevron-right": !isPromptOpen, "fa-chevron-down": isPromptOpen})}></i> Prompt </h4>
    }
    {response.messages.length !==0 && isPromptOpen && <ResponseMessage message={response.messages[0]} />}
    {response.messages.slice(1).map(m => <ResponseMessage key={m.id} message={m}/>)}
    <div className="gui-packed-row">
      <i className="fas fa-user"></i>
      <TextareaAutosize
            className=""
            value={reply}
            onChange={onReplyChange}
            onKeyDown={onReplyKeyDown}
      />
      <button className="small" title="Shift+Enter" onClick={onReplyClicked}>Reply</button>
    </div>
  </div>;
}
