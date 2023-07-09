import { ChangeEvent, useEffect, useState } from "react";
import { Session } from "./data";
import { navigate } from "raviger";
import { loadSession, newResponse, setSessionPrompt, deleteResponse, SessionMonitor } from "./dataloading";
import { ResponseEditor } from "./responseeditor";

export interface Props {
  sessionId: string;
}

let sessionMonitor: SessionMonitor | null = null;

export function SessionEditor({sessionId}: Props): JSX.Element {
  const [session, setSession] = useState<Session | null>(null);
  const [sessionReload, setSessionReload] = useState<number>(0);

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
    sessionMonitor = new SessionMonitor(sessionId, () => {
      sessionReloadCounter++;
      console.log(`Setting sessionReload to ${sessionReloadCounter}`);
      setSessionReload(sessionReloadCounter);
    });
    sessionMonitor.start();
    return () => {
      if (sessionMonitor != null) {
        console.log(`Stopping monitor of ${sessionId}`);
        sessionMonitor.stop();
        sessionMonitor = null;
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

  return <div className="session-editor">
    {session == null && <div>Loading</div>}
    {session && <>
        <div className="session-prompt-pane">
          <h3>Prompt</h3>
          <textarea
            className="char-width-20"
            value={session.prompt}
            onChange={onPromptChange}
          /><br />
          <button className="success" onClick={onSubmitClicked}>Submit</button>
        </div>
        <div className="session-response-pane">
          <h3>Responses</h3>
          {
            session.responses.map(r => <ResponseEditor response={r} key={r.id} onDeleteClicked={onDeleteClicked} />)
          }
        </div>
      </>
    }
  </div>;
}
