import { ChangeEvent, useEffect, useState } from "react";
import { Session } from "./data";
import { navigate } from "raviger";
import { loadSession, setSessionPrompt } from "./dataloading";

export interface Props {
  sessionId: string;
}

export function SessionEditor({sessionId}: Props): JSX.Element {
  const [session, setSession] = useState<Session | null>(null);
  useEffect(() => {
    (async () => {
      const loadedSession = await loadSession(sessionId);
      if (loadedSession == null) {
        navigate("/");
      }
      setSession(loadedSession);
    })();
  }, [sessionId]);

  const onPromptChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    setSession(setSessionPrompt(session as Session, event.target.value));
  }

  return <div className="session-editor">
    {session == null && <div>Loading</div>}
    {session && <>
        <div className="session-prompt-pane">
          <h3>Prompt</h3>
          <textarea
            value={session.prompt}
            onChange={onPromptChange}
          /><br />
          <button className="success">Submit</button>
        </div>
        <div className="session-response-pane">
          <h3>Responses</h3>
        </div>
      </>
    }
  </div>;
}
