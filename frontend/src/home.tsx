import { ModelOverview, SessionOverview } from "./data";
import { SessionEditor } from "./sessioneditor";
import { SessionOverviewList } from "./sessionoverviewlist";
import { NewSessionButton } from "./newsessionbutton";

export interface Props {
  modelOverview: ModelOverview;
  sessionOverview: SessionOverview;
  sessionId: string;
  onSessionChange: ()=> void;
}

export function Home({ modelOverview, sessionOverview, sessionId, onSessionChange }: Props): JSX.Element {
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
          modelOverview={modelOverview}
        />
      </div>
  </div>
  );
}
