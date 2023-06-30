import { useEffect, useState } from "react";
import { Session, SessionOverview } from "./data";
import { SessionEditor } from "./sessioneditor";
import { SessionOverviewList } from "./sessionoverviewlist";
import { loadSession } from "./dataloading";

export interface Props {
  sessionOverview: SessionOverview;
  sessionId: string;
}

export function Home({ sessionOverview, sessionId }: Props): JSX.Element {

  const [session, setSession] = useState<Session | null>(null);
  useEffect(() => {
    (async () => {
      const loadedSession = await loadSession(sessionId);
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
        <button>+ New Session</button>
      </div>
      <div className="session-tab">
        { session && <SessionEditor key={sessionId} session={session} /> }
        { (session ==null) && <span>Loading</span> }
      </div>
  </div>
  );
}
