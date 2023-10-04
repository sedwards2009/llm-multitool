import { Model, ModelOverview, PresetOverview, Response, TemplateOverview, getModelById, isSettingsValid } from "./data";
import { ChangeEvent, KeyboardEventHandler, useState } from "react";
import classNames from "classnames";
import { ResponseMessage } from "./responsemessage";
import TextareaAutosize from "react-textarea-autosize";

export interface Props {
  response: Response;
  modelOverview: ModelOverview;
  presetOverview: PresetOverview;
  templateOverview: TemplateOverview;
  onAbortClicked: () => void;
  onContinueClicked: () => void;
  onDeleteClicked: () => void;
  onReplySubmit: (replyText: string) => void;
}

export function ResponseEditor({response, modelOverview, presetOverview, templateOverview,
    onAbortClicked, onContinueClicked, onDeleteClicked, onReplySubmit: onReply}: Props): JSX.Element {

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

  let isSendEnabled = false;
  let model: Model | null = null;
  let supportsContinue = false;
  if (response.modelSettingsSnapshot != null) {
    const selectedModelId = response.modelSettingsSnapshot.modelId;
    model = getModelById(modelOverview, selectedModelId);

    const selectedPresetId = response.modelSettingsSnapshot.presetId;
    const selectedTemplateId = response.modelSettingsSnapshot.templateId;
    isSendEnabled = (model?.supportsReply === true) && isSettingsValid(modelOverview, presetOverview, templateOverview, selectedModelId,
      selectedPresetId, selectedTemplateId);

    supportsContinue = model?.supportsContinue === true;
  }

  return <div className="card">
    <h3>Response</h3>
    <div className="controls">
      { response.status === "Running" &&
          <>
            <button className="microtool warning" onClick={onAbortClicked}><i className="far fa-stop-circle"></i></button>
            &nbsp;
            <span className="badge success">Running</span>
          </>
      }
      { response.status === "Pending" && <span className="badge warning">Pending</span>}
      { response.status === "Error" && <span className="badge danger">Error</span>}
      { response.status === "Aborted" && <span className="badge warning">Aborted</span>}
      <button className="microtool danger" onClick={onDeleteClicked}><i className="fa fa-times"></i></button>
    </div>
    {response.messages.length !==0 &&
      <h4 className="prompt-header" onClick={onPromptClicked}><i className={classNames({"fa": true, "fa-chevron-right": !isPromptOpen, "fa-chevron-down": isPromptOpen})}></i> Prompt </h4>
    }
    {isPromptOpen && response.modelSettingsSnapshot != null &&
      <div className="gui-layout cols-1-2 response-settings">
        <div>
          <i className="fa fa-robot"></i>&nbsp;&nbsp;Model:
        </div>
        <div>
          {response.modelSettingsSnapshot.modelName}
        </div>

        <div>
          <i className="fa fa-hammer"></i>&nbsp;&nbsp;Task:
        </div>
        <div>
          {response.modelSettingsSnapshot.templateName}
        </div>

        <div>
          <i className="fa fa-tachometer-alt"></i>&nbsp;&nbsp;Creativeness:
        </div>
        <div>
          {response.modelSettingsSnapshot.presetName}
        </div>
      </div>
    }
    {response.messages.length !==0 && isPromptOpen &&
      <ResponseMessage
        message={response.messages[0]}
        onContinueClicked={null}
      />
    }
    {response.messages.slice(1).map((m ,i) =>
      <ResponseMessage
        key={m.id}
        message={m}
        onContinueClicked={supportsContinue && isSendEnabled && response.status === "Done" &&
          response.messages.length-1 === i+1 ? onContinueClicked : null}
      />
    )}
    {isSendEnabled && response.status === "Done" &&
      <div className="response-message">
        <div className="response-message-gutter"><i className="fas fa-user"></i></div>
        <div className="response-message-text gui-packed-row">
          <TextareaAutosize
                className="expand"
                value={reply}
                onChange={onReplyChange}
                onKeyDown={onReplyKeyDown}
          />
          <button className="compact small success" title="Shift+Enter" onClick={onReplyClicked}>Send</button>
        </div>
      </div>
    }
  </div>;
}
