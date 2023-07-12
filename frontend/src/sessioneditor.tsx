import { ChangeEvent, KeyboardEventHandler, useEffect, useState } from "react";
import { Session } from "./data";
import { navigate } from "raviger";
import TextareaAutosize from "react-textarea-autosize";
import { loadSession, newResponse, setSessionPrompt, deleteResponse, SessionMonitor, SessionMonitorState } from "./dataloading";
import { ResponseEditor } from "./responseeditor";

export interface Props {
  sessionId: string;
}

export function SessionEditor({sessionId}: Props): JSX.Element {
  const [session, setSession] = useState<Session | null>(null);
  const [sessionReload, setSessionReload] = useState<number>(0);
  const [sessionMonitor, setSessionMonitor] = useState<SessionMonitor | null>(null);
  const [sessionMonitorState, setSessionMonitorState] = useState<SessionMonitorState>(SessionMonitorState.IDLE);

  const loadSessionData = async () => {
    console.log(`loadSessionData()`);
    const loadedSession = await loadSession(sessionId);
    if (loadedSession == null) {
      navigate("/");
    }
    setSession(loadedSession);
  };

  useEffect(() => {
    loadSessionData();
  }, [sessionId, sessionReload]);

  useEffect(() => {
    let sessionReloadCounter = sessionReload;
    console.log(`Starting to monitor ${sessionId}`);

    const onChange = () => {
      sessionReloadCounter++;
      console.log(`Setting sessionReload to ${sessionReloadCounter}`);
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

  const onDeleteClicked = (responseId: string) => {
    (async () => {
      await deleteResponse((session as Session).id, responseId);
      await loadSessionData();
    })();
  };

  const onSubmitClicked = () => {
    (async () => {
      await newResponse(session as Session);
      await loadSessionData();
    })();
  };

  const onKeyDown: KeyboardEventHandler<HTMLTextAreaElement> = (e) => {
    if (e.code === "Enter" && e.shiftKey) {
      onSubmitClicked();
      e.preventDefault();
    }
  };

  return <div className="session-editor">
    {session == null && <div>Loading</div>}
    {session && <>
        <div className="session-prompt-pane card">
          <h3>Prompt</h3>
          <div className="controls">
            {sessionMonitorState !== SessionMonitorState.CONNECTED &&
              <span className="badge warning">
                <i className="fa fa-plug"></i>
                {" " + SessionMonitorStateToString(sessionMonitorState)}
              </span>
            }
          </div>

          <TextareaAutosize
            className="char-width-20"
            value={session.prompt}
            onChange={onPromptChange}
            onKeyDown={onKeyDown}
          /><br />
          <button className="success" title="Shift+Enter" onClick={onSubmitClicked}>Submit</button>
        </div>
        <div className="session-response-pane">
          {
            session.responses.map(r => <ResponseEditor response={r} key={r.id} onDeleteClicked={onDeleteClicked} />)
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

