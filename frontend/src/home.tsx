import { useEffect, useState } from "react";
import { Session, SessionOverview } from "./data";
import { SessionEditor } from "./sessioneditor";
import { SessionOverviewList } from "./sessionoverviewlist";
import { loadSession } from "./dataloading";
import { navigate } from "raviger";
import { NewSessionButton } from "./newsessionbutton";

export interface Props {
  sessionOverview: SessionOverview;
  sessionId: string;
  onSessionChange: ()=> void;
}

export function Home({ sessionOverview, sessionId, onSessionChange }: Props): JSX.Element {
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

  return (
    <div className="top-layout">
      <div className="session-list">
        <SessionOverviewList
          sessionOverview={sessionOverview}
          selectedSessionId={sessionId}
        />
        <NewSessionButton onSessionChange={onSessionChange} />
      </div>
      <div className="session-tab">
        { session && <SessionEditor key={sessionId} session={session} /> }
        { (session ==null) && <span>Loading</span> }
      </div>
  </div>
  );
}
