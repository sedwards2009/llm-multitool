import ReactMarkdown from "react-markdown";
import { Response } from "./data";
import { useState } from "react";
import classNames from "classnames";

export interface Props {
  response: Response;
  onDeleteClicked: (responseId: string) => void;
}

export function ResponseEditor({response, onDeleteClicked}: Props): JSX.Element {

  const [isPromptOpen, setIsPromptOpen] = useState<boolean>(false);

  const onPromptClicked = () => {
    setIsPromptOpen(!isPromptOpen);
  };

  return <div className="card char-width-20">
    <h3>Response</h3>
    <div className="controls">
      { response.status === "Pending" && <span className="badge warning">Pending</span>}
      { response.status === "Running" && <span className="badge success">Running</span>}
      { response.status === "Error" && <span className="badge danger">Error</span>}
      <button className="microtool danger" onClick={() => onDeleteClicked(response.id)}><i className="fa fa-times"></i></button>
    </div>
    <h4 className="prompt-header" onClick={onPromptClicked}><i className={classNames({"fa": true, "fa-chevron-right": !isPromptOpen, "fa-chevron-down": isPromptOpen})}></i> Prompt </h4>
    {isPromptOpen && <><ReactMarkdown children={response.prompt} /><br /></>}
    <h4>Output</h4>
    <ReactMarkdown children={response.text} />
  </div>;
}
