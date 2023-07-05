import { SessionOverview } from "./data";
import { SessionEditor } from "./sessioneditor";
import { SessionOverviewList } from "./sessionoverviewlist";
import { NewSessionButton } from "./newsessionbutton";

export interface Props {
  sessionOverview: SessionOverview;
  sessionId: string;
  onSessionChange: ()=> void;
}

export function Home({ sessionOverview, sessionId, onSessionChange }: Props): JSX.Element {
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
        <SessionEditor
          key={sessionId}
          sessionId={sessionId}
        />
      </div>
  </div>
  );
}
