import { ChangeEvent, KeyboardEventHandler, useEffect, useState } from "react";
import { ModelOverview, PresetOverview, Session, TemplateOverview, getModelById, isSettingsValid } from "./data";
import { navigate } from "raviger";
import TextareaAutosize from "react-textarea-autosize";
import { loadSession, newResponse, setSessionPrompt, deleteResponse, deleteResponseMessage, SessionMonitor,
  SessionMonitorState, setSessionModel, newMessage, setSessionTemplate, setSessionPreset, continueMessage,
  abortResponse,
  uploadFileToSession,
  deleteSessionAttachedFile} from "./dataloading";
import { ResponseEditor } from "./responseeditor";
import { ModelSettings } from "./modelsettings";
import { FileAttachments } from "./fileattachments";

export interface Props {
  sessionId: string;
  modelOverview: ModelOverview;
  presetOverview: PresetOverview;
  templateOverview: TemplateOverview;
  onSessionDelete: () => void;
  onSessionChange: ()=> void;
}

export function SessionEditor({sessionId, modelOverview, presetOverview, templateOverview, onSessionDelete,
    onSessionChange}: Props): JSX.Element {

  const [session, setSession] = useState<Session | null>(null);
  const [sessionReload, setSessionReload] = useState<number>(0);
  const [_, setSessionMonitor] = useState<SessionMonitor | null>(null);
  const [sessionMonitorState, setSessionMonitorState] = useState<SessionMonitorState>(SessionMonitorState.IDLE);
  const [selectedModelId, setSelectedModelId] = useState<string | null>(null);
  const [selectedTemplateId, setSelectedTemplateId] = useState<string | null>(null);
  const [selectedPresetId, setSelectedPresetId] = useState<string | null>(null);

  const loadSessionData = async () => {
    const loadedSession = await loadSession(sessionId);
    if (loadedSession == null) {
      navigate("/");
    }
    setSession(loadedSession);
    if (loadedSession != null) {
      setSelectedModelId(loadedSession?.modelSettings.modelId);
      setSelectedTemplateId(loadedSession?.modelSettings.templateId);
      setSelectedPresetId(loadedSession?.modelSettings.presetId);
    }
  };

  useEffect(() => {
    loadSessionData();
  }, [sessionId, sessionReload]);

  useEffect(() => {
    let sessionReloadCounter = sessionReload;
    console.log(`Starting to monitor ${sessionId}`);

    const onChange = () => {
      sessionReloadCounter++;
      setSessionReload(sessionReloadCounter);
    };

    const onStateChange = (state: SessionMonitorState) => {
      setSessionMonitorState(state);
    };

    const newSessionMonitor = new SessionMonitor(sessionId, onChange, onStateChange);
    setSessionMonitor(newSessionMonitor);
    newSessionMonitor.start();
    return () => {
      if (newSessionMonitor != null) {
        console.log(`Stopping monitor of ${sessionId}`);
        newSessionMonitor.stop();
        setSessionMonitor(null);
      }
    };
  }, [sessionId]);

  const onPromptChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    setSession(setSessionPrompt(session as Session, event.target.value));
  }

  const onModelChange = (modelId: string) => {
    setSelectedModelId(modelId);
    setSession(setSessionModel(session as Session, modelId));
  };

  const onTemplateChange = (templateId: string) => {
    setSelectedTemplateId(templateId);
    setSession(setSessionTemplate(session as Session, templateId));
  };

  const onAbortClicked = (responseId: string) => {
    (async () => {
      await abortResponse((session as Session).id, responseId);
      await loadSessionData();
    })();
  };

  const onDeleteClicked = (responseId: string) => {
    (async () => {
      await deleteResponse((session as Session).id, responseId);
      await loadSessionData();
    })();
  };

  const onDeleteMessageClicked = (responseId: string, messageId: string) => {
    (async () => {
      await deleteResponseMessage((session as Session).id, responseId, messageId);
      await loadSessionData();
    })();
  };

  const onSubmitClicked = () => {
    (async () => {
      await newResponse(session as Session);
      await loadSessionData();
      onSessionChange();
    })();
  };

  const onKeyDown: KeyboardEventHandler<HTMLTextAreaElement> = (e) => {
    if (e.code === "Enter" && e.shiftKey) {
      onSubmitClicked();
      e.preventDefault();
    }
  };

  const onReplySubmit = (responseId: string, reply: string) => {
    (async () => {
      await newMessage(session as Session, responseId, reply);
      await loadSessionData();
    })();
  };

  const onContinueClicked = (responseId: string) => {
    (async () => {
      await continueMessage(session as Session, responseId);
      await loadSessionData();
    })();
  };

  const onUploadFile = (file: File) => {
    (async () => {
      await uploadFileToSession(session as Session, file);
      await loadSessionData();
    })();
  };

  const onDeleteFile = (filename: string) => {
    (async () => {
      await deleteSessionAttachedFile((session as Session).id, filename);
      await loadSessionData();
    })();
  };

  const onPresetChange = (presetId: string) => {
    setSelectedPresetId(presetId);
    setSession(setSessionPreset(session as Session, presetId));
  };

  const isSendEnabled = isSettingsValid(modelOverview, presetOverview, templateOverview, selectedModelId,
    selectedPresetId, selectedTemplateId);

  const model = getModelById(modelOverview, selectedModelId);
  const supportImages = model?.supportsImages == true;

  return <div className="session-editor">
    {session == null && <div>Loading</div>}
    {session && <>
        <div className="session-prompt-pane card">
          <div className="controls">
            <button className="microtool danger" onClick={onSessionDelete}><i className="fa fa-times"></i></button>
          </div>

          <h3>Settings</h3>
          <ModelSettings
            modelOverview={modelOverview}
            templateOverview={templateOverview}
            presetOverview={presetOverview}
            selectedModelId={selectedModelId}
            setSelectedModelId={onModelChange}
            selectedTemplateId={selectedTemplateId}
            setSelectedTemplateId={onTemplateChange}
            selectedPresetId={selectedPresetId}
            setSelectedPresetId={onPresetChange}
          />
          <h3>Prompt</h3>
          <div className="controls">
            {sessionMonitorState !== SessionMonitorState.CONNECTED &&
              <span className="badge warning">
                <i className="fa fa-plug"></i>
                {" " + SessionMonitorStateToString(sessionMonitorState)}
              </span>
            }
          </div>

          <div className="gui-packed-row">
            <TextareaAutosize
              className="expand"
              value={session.prompt}
              onChange={onPromptChange}
              onKeyDown={onKeyDown}
            />
            <button
              className="compact small success"
              title="Shift+Enter"
              onClick={onSubmitClicked}
              disabled={!isSendEnabled}>
                Send
            </button>
          </div>
          {supportImages && <FileAttachments
            sessionId={sessionId}
            onUploadFile={onUploadFile}
            onDeleteFile={onDeleteFile}
            attachedFiles={session.attachedFiles}
          />}
        </div>
        <div className="session-response-pane">
          {
            [...session.responses].reverse().map(r => <ResponseEditor
              sessionId={sessionId}
              response={r}
              key={r.id}
              modelOverview={modelOverview}
              presetOverview={presetOverview}
              templateOverview={templateOverview}
              onAbortClicked={() => onAbortClicked(r.id)}
              onDeleteClicked={() => onDeleteClicked(r.id)}
              onDeleteMessageClicked={(messageId: string) => onDeleteMessageClicked(r.id, messageId)}
              onReplySubmit={text => onReplySubmit(r.id, text)}
              onContinueClicked={() => onContinueClicked(r.id)}
            />)
          }
        </div>
      </>
    }
  </div>;
}

function SessionMonitorStateToString(state: SessionMonitorState): string {
  return {
    [SessionMonitorState.IDLE]: "Idle",
    [SessionMonitorState.CONNECTING]: "Connecting",
    [SessionMonitorState.CONNECTED]: "Connected",
    [SessionMonitorState.WAITING_TO_RECONNECT]: "Waiting to reconnect",
  }[state];
}

